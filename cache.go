package batch

import (
	"bufio"
	"errors"
	"os"
)

func populateProcessedRecordsCache(file string, processedRecordsCache *stringSet) {
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		LogPrintFunc("skipping processed records cache population as no known previously processed records")
		close(processedRecordsCache.Ready)
	} else {
		err = os.Chmod(file, 0666)
		check(err, LogFatalFunc, "fatal error setting file permissions on %v: %v\n", file, err)

		f, err := os.Open(file)

		defer func() {
			f.Close()
			LogPrintFunc("processed records cache ready")
			close(processedRecordsCache.Ready)
		}()

		check(err, LogFatalFunc, "fatal error opening %v for read: %v\n", file, err)

		s := bufio.NewScanner(f)

		LogPrintFunc("populating processed records cache")

		for s.Scan() {
			if s.Text() != "" {
				processedRecordsCache.add(s.Text())
			}
		}

		check(s.Err(), LogFatalFunc, "fatal error parsing data from %v to populate processed records cache: %v\n", file, s.Err())

		LogPrintFunc("populated processed records cache with %v records", processedRecordsCache.size())
	}
}
