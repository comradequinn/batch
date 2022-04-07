package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/comradequinn/batch"
)

func main() {
	// numbers represents two numbers to be added
	type numbers struct {
		Val1 int
		Val2 int
	}

	cfg := batch.Config[numbers]{
		Workers:                 10,
		MinRecordProcessingTime: time.Second,
		InputFile:               "./data/numbers.csv",
		ProcessedRecordKeysFile: "./data/done.dat",
		KeyFor: func(n numbers) (string, error) {
			return fmt.Sprintf("%v:%v\n", n.Val1, n.Val2), nil
		},
		Parse: func(line []string) (numbers, error) {
			val1, _ := strconv.Atoi(line[0])
			val2, _ := strconv.Atoi(line[1])

			return numbers{Val1: val1, Val2: val2}, nil
		},
		Task: func(n numbers) error {
			log.Printf("executing addition task: %v + %v = %v", n.Val1, n.Val2, n.Val1+n.Val2)

			return nil
		},
	}

	batch.Run(cfg)
}
