package batch

import (
	"bufio"
	"os"
	"time"
)

// startProcessedWorker returns a channel which marks records sent to it as processed by writing their identifying data to file
func startProcessedWorker[T any](buffer int, file string, progressReportFrequency time.Duration, started time.Time, processedRecordsCache *stringSet, keyFor func(T) (string, error)) (chan<- T, <-chan struct{}) {
	in, done := make(chan T, buffer), make(chan struct{}, 1)

	populateProcessedRecordsCache(file, processedRecordsCache)

	if progressReportFrequency == 0 {
		progressReportFrequency = time.Second
	}

	progressReportDue, recordsProcessed := time.NewTicker(progressReportFrequency), int64(0)

	go func() {
		<-processedRecordsCache.Ready

		LogPrintFunc("starting processed records process")

		f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModeAppend)
		check(err, LogFatalFunc, "fatal error opening %v for write: %v\n", file, err)

		defer f.Close()

		err = os.Chmod(file, 0666)
		check(err, LogFatalFunc, "fatal error setting file permissions on %v: %v\n", file, err)

		w := bufio.NewWriter(f)

		for record := range in {
			recordsProcessed++

			key, err := keyFor(record)
			check(err, LogFatalFunc, "fatal error creating key for %+v from file %v: %v\n", record, file, err)

			_, err = w.WriteString(key + "\n")
			check(err, LogFatalFunc, "fatal error writing to %v: %v\n", file, err)

			err = w.Flush()
			check(err, LogFatalFunc, "fatal error flushing to %v: %v\n", file, err)

			select {
			case <-progressReportDue.C:
				go func(rp int64, since time.Duration) {
					ReportProgressFunc("processed %v records in %v", rp, since)
				}(recordsProcessed, time.Since(started))
			default:
			}
		}

		done <- struct{}{}

		LogPrintFunc("processed records process ended")
	}()

	return in, done
}
