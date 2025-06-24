package csv

import "encoding/csv"

// Options is a struct to hold optional parameters for CommonCSV initialization
type Options struct {
	Headers          Headers
	BatchSize        int
	CSVReader        *csv.Reader
	fileInfo         FileInfo
	localFilePath    string
	source           Source
	headerSanitizers []sanitizerFunc
	rowSantiziers    []sanitizerFunc
	rawData          []byte
	delimiter        rune
	rawURLInfo       RawURLInfo
}

// OptionFunc is a function type for setting options
type OptionFunc func(*Options)

// WithHeaders sets the headers option
func WithHeaders(headers Headers) OptionFunc {
	return func(o *Options) {
		o.Headers = headers
	}
}

// WithBatchSize sets the batch size option
func WithBatchSize(batchSize int) OptionFunc {
	return func(o *Options) {
		o.BatchSize = batchSize
	}
}

// WithCSVReader sets the csv.Reader option
func WithCSVReader(reader *csv.Reader) OptionFunc {
	return func(o *Options) {
		o.CSVReader = reader
	}
}

func WithFileInfo(objectKey, bucket string) OptionFunc {
	return func(o *Options) {
		o.fileInfo = FileInfo{
			ObjectKey: objectKey,
			Bucket:    bucket,
		}
	}
}

func WithLocalFileInfo(filePath string) OptionFunc {
	return func(o *Options) {
		o.localFilePath = filePath
	}
}

func WithSource(source Source) OptionFunc {
	return func(o *Options) {
		o.source = source
	}
}

func WithRawData(rawData []byte) OptionFunc {
	return func(o *Options) {
		o.rawData = rawData
	}
}

func WithRawURLInfo(url, encoding string) OptionFunc {
	return func(o *Options) {
		o.rawURLInfo.URL = url
		o.rawURLInfo.Encoding = encoding
	}
}

func WithDelimiter(delimiter rune) OptionFunc {
	return func(o *Options) {
		o.delimiter = delimiter
	}
}

func WithHeaderSanitizers(sanitizers ...sanitizerFunc) OptionFunc {
	return func(o *Options) {
		o.headerSanitizers = append(o.headerSanitizers, sanitizers...)
	}
}

func WithDataRowSanitizers(sanitizers ...sanitizerFunc) OptionFunc {
	return func(o *Options) {
		o.rowSantiziers = append(o.rowSantiziers, sanitizers...)
	}
}
