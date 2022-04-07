# Overview

batch provides generalised batch record processing functionality:

* Lines are read from a specified file and passed to the specified parse function to return a record type
* Each record is then passed to a specified task function
  * The task functions are invoked with the specified levels of concurrency and/or delay
* When the task has completed, a key is derived from the record using the specified KeyFor function
* The key is written to a specified file used to record processed records
* If the batch is re-run (without removing the processed records file), any keys already present in the processed records file will not be re-processed, allowing the job to stopped, started or recover from a termination

# Examples

The below code illustrates how to consume the library. This can be found in full in the `./cmd/example` directory and executed with `make example`

```go
package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/comradequinn/batch"
)

// the simple example below reads in lines of number pairs from a file and adds each pair together, logging the result to stdout
func main() {
  // numbers represents two numbers to be added
  type numbers struct {
	  Val1 int
	  Val2 int
  }

	cfg := batch.Config{
		Workers:                 10, // specifies the number of workers processng the record
		MinRecordProcessingTime: time.Second,  // specifies the minimum execution time required for each record a worker processes
		InputFile:               "./data/numbers.csv", // this contains series of lines numbers to be added in the format "1,2"
		ProcessedRecordKeysFile: "./data/done.dat", // this is where the keys of processed number records will be written so the job can be stopped and restarted safely
		KeyFor:                  keyFor, // this func is used to derive a key from a number type
		Parse:                   parse, // this func is used to parse a record from 
		Task:                    task, // this func performs the addition of the numnbers and the logging
	}

	batch.Run(cfg)
}

// Below are the implementations for KeyFor, Parse & Task

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

```



