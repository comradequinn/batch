package batch

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func Test_Run(t *testing.T) {
	type numbers struct {
		Val1 int
		Val2 int
	}

	inf, outf := "./testdata/numbers.csv", "./testdata/done.dat"
	expectedTotal, actualTotal, mx := 203, 0, sync.Mutex{}

	os.Remove(outf) // clear any cached progress

	cfg := Config[numbers]{
		Workers:                 10,
		MinRecordProcessingTime: time.Millisecond * 100,
		InputFile:               inf,
		ProcessedRecordKeysFile: outf,
		Parse: func(line []string) (numbers, error) {
			val1, _ := strconv.Atoi(line[0])
			val2, _ := strconv.Atoi(line[1])

			return numbers{Val1: val1, Val2: val2}, nil
		},
		KeyFor: func(n numbers) (string, error) {
			return fmt.Sprintf("%v:%v\n", n.Val1, n.Val2), nil
		},
		Task: func() func(numbers) error { // this task adds all number presented to it and stores them in `actualTotal`
			return func(n numbers) error { // if `actualTotal` matches `expectedTotal` then all records were read, parsed correctly and passed to the task
				mx.Lock()
				defer mx.Unlock()

				actualTotal += n.Val1 + n.Val2

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
