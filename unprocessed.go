package batch

import (
	"bufio"
	"os"
	"strings"
)

// startUnprocessedWorker returns a channel from which unprocessed records can be read
func startUnprocessedWorker[T any](bufferSize int, file, fileDelimiter string, continueOnError bool, processedRecordsCache *stringSet, keyFor func(T) (string, error), parse func([]string) (T, error)) <-chan T {
	<-processedRecordsCache.Ready

	out := make(chan T, bufferSize)

	go func() {
		LogPrintFunc("starting read unprocessed records process")

		defer func() {
			err := recover()
			check(err, LogFatalFunc, "fatal error reading from %v: %v\n", file, err)
		}()

		f, err := os.Open(file)
		check(err, LogFatalFunc, "fatal error opening %v: %v\n", file, err)
		defer f.Close()

		s := bufio.NewScanner(f)

		if fileDelimiter == "" {
			fileDelimiter = ","
		}

		for s.Scan() {
			line := strings.Split(strings.TrimPrefix(s.Text(), string('\ufeff')), fileDelimiter) // strip bom if present

			record, err := parse(line)

			if err != nil {
				LogPrintFunc("error parsing line %v: %v\n", line, err)

				if !continueOnError {
					os.Exit(1)
				}

				continue
			}

			key, err := keyFor(record)
			check(err, LogFatalFunc, "fatal error creating key for %+v from file %v: %v\n", record, file, err)

			if !processedRecordsCache.has(key) {
				out <- record
			}
		}

		check(s.Err(), LogFatalFunc, "fatal error parsing data from %v: %v\n", file, s.Err())

		close(out)

		LogPrintFunc("read unprocessed records process ended")
	}()

	return out
}
