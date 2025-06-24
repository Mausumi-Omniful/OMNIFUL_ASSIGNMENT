package csv

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/log"
	"strings"
)

type DestinationType string
type Destination struct {
	fileName           string
	bucket             string
	uploadDirectory    string
	destType           DestinationType
	randomizedFileName bool
}

const (
	DestinationS3    DestinationType = "s3"
	DestinationLocal DestinationType = "local"
)

// String returns the string representation of the Source type.
func (destination *Destination) String() string {
	return string(destination.destType)
}

// IsDestinationLocal checks if the Source is set to "local".
func IsDestinationLocal(destination Destination) bool {
	return destination.destType == DestinationLocal
}

// IsDestinationS3 checks if the Source is set to "s3".
func IsDestinationS3(destination Destination) bool {
	return destination.destType == DestinationS3
}

var DestinationMap = map[string]DestinationType{
	"s3":    DestinationS3,
	"local": DestinationS3,
}

// GetDestination converts a string to the Source type. By default we are assuming local as destination
func GetDestination(s string) DestinationType {
	destination, found := DestinationMap[strings.ToLower(s)]
	if !found {
		return DestinationLocal
	}

	return destination
}

func (destination *Destination) GetFileName() string {
	return destination.fileName
}

func (destination *Destination) GetBucket() string {
	return destination.bucket
}

func (destination *Destination) GetType() DestinationType {
	return destination.destType
}

func (destination *Destination) IsRandomizedFileName() bool {
	return destination.randomizedFileName
}

func (destination *Destination) GetUploadDirectory() string {
	return destination.uploadDirectory
}

func (destination *Destination) SetFileName(filename string) (err error) {
	if !isValidCSVName(filename) {
		log.Errorf("invalid file Name %s", filename)
		return fmt.Errorf("invalid file Name %s", filename)
	}
	destination.fileName = filename

	return nil
}

func (destination *Destination) SetBucket(bucket string) {
	destination.bucket = bucket
}

func (destination *Destination) SetUploadDirectory(directory string) {
	destination.uploadDirectory = addTrailingSlash(directory)
}

func (destination *Destination) SetRandomizedFileName(set bool) {
	destination.randomizedFileName = set
}

func isValidCSVName(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".csv")
}

func generateUniqueFilename() string {
	filename := fmt.Sprintf("%s.csv", uuid.New().String())
	return filename
}

type DestinationOptions struct {
	fileName           string
	bucket             string
	directory          string
	destType           DestinationType
	randomizedFileName bool
}

type DestinationOptionFunc func(*DestinationOptions)

func WithFileName(fileName string) DestinationOptionFunc {
	return func(o *DestinationOptions) {
		o.fileName = fileName
	}
}

func WithBucketName(bucket string) DestinationOptionFunc {
	return func(o *DestinationOptions) {
		o.bucket = bucket
	}
}

func WithUploadDirectory(directory string) DestinationOptionFunc {
	return func(o *DestinationOptions) {
		o.directory = directory
	}
}

func WithType(destType DestinationType) DestinationOptionFunc {
	return func(o *DestinationOptions) {
		o.destType = destType
	}
}

func WithRandomizedFileName(set bool) DestinationOptionFunc {
	return func(o *DestinationOptions) {
		o.randomizedFileName = set
	}
}

func NewDestination(options ...DestinationOptionFunc) (*Destination, error) {
	opts := &DestinationOptions{}

	for _, option := range options {
		option(opts)
	}

	if opts.randomizedFileName {
		opts.fileName = generateUniqueFilename()
	}

	if !isValidCSVName(opts.fileName) {
		log.Errorf("invalid file Name %s", opts.fileName)
		return nil, fmt.Errorf("invalid file Name %s", opts.fileName)
	}

	return &Destination{
		fileName:           opts.fileName,
		bucket:             opts.bucket,
		destType:           opts.destType,
		randomizedFileName: opts.randomizedFileName,
		uploadDirectory:    addTrailingSlash(opts.directory),
	}, nil
}
