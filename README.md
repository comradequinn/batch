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

	cfg := batch.Config[numbers]{
		Workers:                 10, // specifies the number of workers required to be concurrently processing records
		MinRecordProcessingTime: time.Second, // specifies the minimum time to be spent on each record by a worker
		InputFile:               "./data/numbers.csv", // this csv file contains a series of lines, each containing two numbers, eg: "1,2"
		ProcessedRecordKeysFile: "./data/done.dat", // this is where the keys of processed records will be written so the job can be stopped and restarted safely
		KeyFor: func(n numbers) (string, error) { // this func is used to derive a key from a numbers type
			return fmt.Sprintf("%v:%v\n", n.Val1, n.Val2), nil
		},
		Parse: func(line []string) (numbers, error) { // this func is used to parse a numbers type from a line in the input file
			val1, _ := strconv.Atoi(line[0]) 
			val2, _ := strconv.Atoi(line[1]) // ignore errors for the purposes of the example

			return numbers{Val1: val1, Val2: val2}, nil
		},
		Task: func(n numbers) error { // this func performs the addition of the two numnbers in a numbers type and logs the result
			log.Printf("executing addition task: %v + %v = %v", n.Val1, n.Val2, n.Val1+n.Val2)

			return nil
		},
	}

	batch.Run(cfg) // this starts the batch job and blocks until it is completed
}

```



