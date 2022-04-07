package batch

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

// numbers represents two numbers to be added
type numbers struct {
	Val1 int
	Val2 int
}

func Test_Run(t *testing.T) {
	inf, outf := "./testdata/numbers.csv", "./testdata/done.dat"
	expectedTotal, actualTotal, mx := 203, 0, sync.Mutex{}

	os.Remove(outf) // clear any cached progress

	cfg := Config{
		Workers:                 10,
		MinRecordProcessingTime: time.Millisecond * 100,
		InputFile:               inf,
		ProcessedRecordKeysFile: outf,
		Parse: func(line []string) (interface{}, error) {
			n := numbers{}

			n.Val1, _ = strconv.Atoi(line[0])
			n.Val2, _ = strconv.Atoi(line[1])

			return n, nil
		},
		KeyFor: func(n interface{}) (string, error) {
			numbers, _ := n.(numbers)
			return fmt.Sprintf("%v:%v\n", numbers.Val1, numbers.Val2), nil
		},
		Task: func() func(interface{}) error { // this task adds all number presented to it and stores them in `actualTotal`
			return func(n interface{}) error { // if `actualTotal` matches `expectedTotal` then all records were read, parsed correctly and passed to the task
				numbers, _ := n.(numbers)

				mx.Lock()
				defer mx.Unlock()

				actualTotal += numbers.Val1 + numbers.Val2

				return nil
			}
		}(),
	}

	Run(cfg)

	if expectedTotal != actualTotal {
		t.Fatalf("the expected total of the processed numbers was not correct on a clean batch run. expected %v. got %v", expectedTotal, actualTotal)
	}

	Run(cfg) // run again, but this time the cached progress should mean no records are processed

	expectedTotal, actualTotal = 0, 0

	if expectedTotal != actualTotal {
		t.Fatalf("the expected total of the processed numbers was not correct on a previously executed batch run. expected %v. got %v", expectedTotal, actualTotal)
	}

	os.Remove(outf) // clear any cached progress
}
