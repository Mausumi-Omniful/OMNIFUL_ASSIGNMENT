package csv

import (
	"encoding/csv"
	"os"
)

// DataProvider is the interface for a CSV data provider.
type DataProvider interface {
	// GetHeaderRow returns the header row for the CSV.
	GetHeaderRow() []string

	// GetDataRows returns the data rows for the CSV.
	GetDataRows() [][]string
}

// WriterOptions is a struct to hold optional parameters for CommonCSV initialization
type WriterOptions struct {
	Writer          *csv.Writer
	Headers         Headers
	file            *os.File
	destination     Destination
	isHeaderWritten bool
	shouldAddBOM    bool
}

// WriterOptionFunc is a function type for setting WriterOptions
type WriterOptionFunc func(*WriterOptions)

// WithWriterHeaders sets the headers option
func WithWriterHeaders(headers Headers) WriterOptionFunc {
	return func(o *WriterOptions) {
		o.Headers = headers
	}
}

func WithWriterDestination(destination Destination) WriterOptionFunc {
	return func(o *WriterOptions) {
		o.destination = destination
	}
}

// WithBOMCharacter configures the Writer to optionally include a BOM at the beginning of the output file.
func WithBOMCharacter(shouldAddBOM bool) WriterOptionFunc {
	return func(o *WriterOptions) {
		o.shouldAddBOM = shouldAddBOM
	}
}
