package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/comradequinn/batch"
)

// numbers represents two numbers to be added
type numbers struct {
	Val1 int
	Val2 int
}

func main() {
	cfg := batch.Config{
		Workers:                 10,
		MinRecordProcessingTime: time.Second,
		InputFile:               "./data/numbers.csv",
		ProcessedRecordKeysFile: "./data/done.dat",
		KeyFor:                  keyFor,
		Parse:                   parse,
		Task:                    task,
	}

	batch.Run(cfg)
}

func task(n interface{}) error {
	numbers, ok := n.(numbers)

	if !ok {
		return fmt.Errorf("unable to convert %+v to type numbers", n)
	}

	log.Printf("executing addition task: %v + %v = %v", numbers.Val1, numbers.Val2, numbers.Val1+numbers.Val2)

	return nil
}

func keyFor(n interface{}) (string, error) {
	numbers, ok := n.(numbers)

	if !ok {
		return "", fmt.Errorf("unable to convert %+v to type numbers", n)
	}

	return fmt.Sprintf("%v:%v\n", numbers.Val1, numbers.Val2), nil
}

func parse(line []string) (interface{}, error) {
	var err error

	n := numbers{}

	n.Val1, err = strconv.Atoi(line[0])

	if err != nil {
		return nil, fmt.Errorf("fatal error reading %v as val1: %v", line[0], err)
	}

	n.Val2, err = strconv.Atoi(line[1])

	if err != nil {
		return nil, fmt.Errorf("fatal error reading %v as val2: %v", line[1], err)
	}

	return n, nil
}
