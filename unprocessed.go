package batch

import (
	"bufio"
	"os"
	"strings"
)

// startUnprocessedWorker returns a channel from which unprocessed records can be read
func startUnprocessedWorker(bufferSize int, file, fileDelimiter string, processedRecordsCache *stringSet, keyFor func(interface{}) (string, error), parse func([]string) (interface{}, error)) <-chan interface{} {
	<-processedRecordsCache.Ready

	out := make(chan interface{}, bufferSize)

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

			check(err, LogFatalFunc, "fatal error opening %v: %v\n", file, err)

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
