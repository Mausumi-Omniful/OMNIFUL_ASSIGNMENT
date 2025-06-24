package csv

import (
	"strings"
)

const (
	SourceS3       = "s3"
	SourceLocal    = "local"
	SourceS3Large  = "s3-large"
	SourceRawBytes = "raw-bytes"
	SourceRawURL   = "raw-url"
)

type Source string

const (
	S3       Source = SourceS3
	Local    Source = SourceLocal
	S3Large  Source = SourceS3Large
	RawBytes Source = SourceRawBytes
	RawURL   Source = SourceRawURL
)

// String returns the string representation of the Source type.
func (s Source) String() string {
	return string(s)
}

// IsSourceLocal checks if the Source is set to "local".
func IsSourceLocal(source Source) bool {
	return source == Local
}

// IsSourceS3 checks if the Source is set to "s3".
func IsSourceS3(source Source) bool {
	return source == S3
}

var sourceMap = map[string]Source{
	SourceS3:       S3,
	SourceLocal:    Local,
	SourceS3Large:  S3Large,
	SourceRawBytes: RawBytes,
}

// GetSource converts a string to the Source type. By default we are assuming s3 as source
func GetSource(s string) Source {
	source, found := sourceMap[strings.ToLower(s)]
	if !found {
		return S3
	}

	return source
}
