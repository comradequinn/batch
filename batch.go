// Package batch provides generalised batch record processing functionality.
// - Lines are read from a client provided file and passed to the client provided parse function to return a record type
// - Each record is then passed to a client provided task function
// 		- The task functions are invoked with the client provided levels of concurrency and delay
// - When the task has completed, a key is derived from the record using the client provided KeyFor function
// - The key is written to a client provided file used to record processed records
// - When the batch is run, any keys already present in the processed records file will not be re-processed, allowing the job to stopped, started or recover from a termination
package batch

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	// LogFatalFunc is the func invoked when a fatal error occurs. It must terminate the process
	LogFatalFunc = log.Fatalf
	// LogPrintFunc is the func invoked when relevant information is generated occurs. It must not terminate the process
	LogPrintFunc = log.Printf
	// ReportProgressFunc is the func invoked periodically to report on batch progress
	ReportProgressFunc = log.Printf
)

func Run(cfg Config) {
	LogPrintFunc("starting batch record processing from %v\n", cfg.InputFile)

	// wrap KeyFor func to ensure it trims whitespace
	cfg.KeyFor = func(original func(interface{}) (string, error)) func(interface{}) (string, error) {
		return func(i interface{}) (string, error) {
			key, err := original(i)
			return strings.TrimSpace(key), err
		}
	}(cfg.KeyFor)

	processedRecordsCache, bufferSize, wg := newStringSet(), cfg.Workers*2, sync.WaitGroup{}

	processedRecords, done := startProcessedWorker(bufferSize, cfg.ProcessedRecordKeysFile, cfg.ProgressReportFrequency, time.Now(), processedRecordsCache, cfg.KeyFor)
	unprocessedRecords := startUnprocessedWorker(bufferSize, cfg.InputFile, processedRecordsCache, cfg.KeyFor, cfg.Parse)

	for i := 0; i < cfg.Workers; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			initialMinRecordProcessingTime := cfg.MinRecordProcessingTime

			if initialMinRecordProcessingTime == 0 {
				initialMinRecordProcessingTime = time.Second
			}

			// avoid all routines repeatedly executing their task at same time by blocking each routine for a different period before it starts
			staggeredStartDuration, tickerReset := time.Duration(int64(initialMinRecordProcessingTime)/int64(cfg.Workers)*int64(workerID)), false
			minBatchItemTimeElapsed, recordsProcessed := time.NewTicker(staggeredStartDuration), 0

			LogPrintFunc("starting worker %v after initial dalay of %vms", workerID, staggeredStartDuration.Milliseconds())

			for record := range unprocessedRecords {
				if err := cfg.Task(record); err != nil {
					LogPrintFunc("error executing task for batch record %+v on worker %v\n", record, workerID)

					if !cfg.ContinueOnError {
						os.Exit(1)
					}
				}

				processedRecords <- record
				recordsProcessed++

				if cfg.MinRecordProcessingTime > 0 {
					<-minBatchItemTimeElapsed.C // if the minimum record processing time has not yet elapsed, yield until it does

					if !tickerReset {
						minBatchItemTimeElapsed.Reset(cfg.MinRecordProcessingTime)
						tickerReset = true
					}
				}
			}

			LogPrintFunc("worker %v stopped after processing %v records", workerID, recordsProcessed)
		}(i + 1)
	}

	wg.Wait()
	close(processedRecords)

	<-done

	LogPrintFunc("batch record processing completed\n")
}
