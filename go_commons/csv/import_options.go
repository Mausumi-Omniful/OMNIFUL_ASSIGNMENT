package csv

// ImporterOptions is a struct to hold optional parameters for CommonCSV initialization
type ImporterOptions struct {
	oReader             *CommonCSV
	oWriter             *CommonCSVWriter
	notifyFailedEntries bool
}
type ImporterOption func(*ImporterOptions)

func WithReader(reader *CommonCSV) ImporterOption {
	return func(i *ImporterOptions) {
		i.oReader = reader
	}
}

func WithWriter(writer *CommonCSVWriter) ImporterOption {
	return func(i *ImporterOptions) {
		i.oWriter = writer
	}
}

func WithNotifyFailedRows(notifyFailedEntries bool) ImporterOption {
	return func(i *ImporterOptions) {
		i.notifyFailedEntries = notifyFailedEntries
	}
}
