package batch

import (
	"time"
)

// Config defines the configuration options for the batch
type Config[T any] struct {
	// Workers states the number of workers available to process records in a batch
	Workers int
	// MinRecordProcessingTime states the minimum time each worker should spend on processing a single record.
	// If the record processing does not organically take this duration, it will block until it has elapsed
	MinRecordProcessingTime time.Duration
	// ProgressReportFrequency states the frequency with which progess reports will be written
	ProgressReportFrequency time.Duration
	// InputFile contains the records that are to be batch processed
	InputFile string
	// InputFileDelimiter specifies the delimiter of fields within each line of the input file, default is ","
	InputFileDelimiter string
	// ProcessedRecordKeysFile will contain unique keys for each record that has been processed
	ProcessedRecordKeysFile string
	// KeyFor represents the function to a derive a unique key from a record
	KeyFor func(T) (string, error)
	// Parse represents the function to parse a record into the type required by the task
	Parse func([]string) (T, error)
	// Task represents the function to which each parsed record is passed in order to be acted upon
	Task func(T) error
	// ContinueOnError states whether or not processing will terminate if a Task invocation or Parse attempt returns an error. If false, a log will still be written
	ContinueOnError bool
}
