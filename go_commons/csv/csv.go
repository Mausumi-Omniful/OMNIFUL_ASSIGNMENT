package csv

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/omniful/go_commons/config"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/omniful/go_commons/log"
)

const defaultBatchSize = 100
const defaultDelimiter = ','

// Headers represents a slice of strings for CSV headers
type Headers []string

// Records represents a 2D slice of strings for CSV records
type Records [][]string

// ToMaps converts records to a slice of maps using headers as keys
func (records Records) ToMaps(headers Headers) []map[string]string {
	var keyValueMaps []map[string]string

	for _, record := range records {
		keyValueMap := make(map[string]string)

		// Iterate through the header and record values to create key-value pairs
		for i, headerValue := range headers {
			keyValueMap[headerValue] = record[i]
		}

		// Append the key-value map to the slice
		keyValueMaps = append(keyValueMaps, keyValueMap)
	}

	return keyValueMaps
}

func (records Records) Unmarshal(headers Headers, data interface{}) (err error) {
	if records == nil {
		err = fmt.Errorf("nil pointer reference")
		return
	}

	recordMap := records.ToMaps(headers)
	// Convert the JSON objects to a JSON-encoded string
	jsonData, err := json.Marshal(recordMap)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	err = json.Unmarshal(jsonData, data)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

func containsBOM(s string) bool {
	if len(s) < 3 {
		return false
	}

	return s[0] == '\xEF' && s[1] == '\xBB' && s[2] == '\xBF'
}

func removeBOM(s string) string {
	if containsBOM(s) {
		return s[3:]
	}
	return s
}

func (headers *Headers) SanitizeSpace() {
	if headers == nil {
		return
	}

	for i, header := range *headers {
		if containsBOM(header) {
			header = removeBOM(header)
		}

		(*headers)[i] = strings.TrimSpace(header)
	}
}

type FileInfo struct {
	ObjectKey string
	Bucket    string
}

type CommonCSV struct {
	Reader           *csv.Reader
	Headers          Headers
	batchSize        int
	eof              bool
	rawData          []byte
	fileInfo         FileInfo
	localFilePath    string
	file             *os.File
	source           Source
	headerSanitizers []sanitizerFunc
	rowSanitizers    []sanitizerFunc
	delimiter        rune
	rawURLInfo       RawURLInfo
}

// RawURLInfo holds metadata for CSV files obtained from a raw URL source.
// It is used within the package to manage CSV files that are fetched over HTTP.
// The struct stores the URL from which the CSV file is downloaded, the encoding type (for example "gzip"),
// and local file paths for the downloaded and the optionally decoded files.
type RawURLInfo struct {
	URL                string // URL to fetch the CSV file from.
	Encoding           string // Encoding applied to the CSV file (e.g., "gzip").
	decodedFilePath    string // Path to the decoded CSV file after processing.
	downloadedFilePath string // Path to the temporarily stored downloaded CSV file.
}

func (commonCSV *CommonCSV) GetSource() Source {
	return commonCSV.source
}

// NewCommonCSV initializes a CommonCSV instance with options
func NewCommonCSV(options ...OptionFunc) (*CommonCSV, error) {
	opts := &Options{
		BatchSize:        defaultBatchSize, // Default value
		source:           SourceS3,         // Default Source
		headerSanitizers: []sanitizerFunc{sanitizeBOM, SanitizeSpace},
		delimiter:        defaultDelimiter,
	}

	for _, option := range options {
		option(opts)
	}

	// Use the provided csv.Reader if available, otherwise create a new one
	csvReader := opts.CSVReader
	if opts.Headers != nil {
		csvReader.FieldsPerRecord = len(opts.Headers)
	}

	return &CommonCSV{
		Reader:           opts.CSVReader,
		Headers:          opts.Headers,
		batchSize:        opts.BatchSize,
		fileInfo:         opts.fileInfo,
		localFilePath:    opts.localFilePath,
		source:           opts.source,
		headerSanitizers: opts.headerSanitizers,
		rowSanitizers:    opts.rowSantiziers,
		rawData:          opts.rawData,
		delimiter:        opts.delimiter,
		rawURLInfo:       opts.rawURLInfo,
	}, nil
}

// ReadNextBatch reads the next batch of records from the CSV data
func (commonCSV *CommonCSV) ReadNextBatch() (records Records, err error) {
	records = make(Records, 0)
	if commonCSV == nil {
		return
	}

	if commonCSV.Headers == nil {
		commonCSV.ParseHeaders()
	}

	for i := 0; i < commonCSV.batchSize; i++ {
		record, err := commonCSV.Reader.Read()
		if err == io.EOF {
			commonCSV.SetEOF()
			break
		} else if err != nil {
			return nil, fmt.Errorf("error reading CSV record: %w", err)
		}

		record = Sanitize(record, commonCSV.rowSanitizers)
		records = append(records, record)
	}

	return records, nil
}

func (commonCSV *CommonCSV) SetEOF() {
	if commonCSV == nil {
		return
	}

	// Closing file in-case csv is being read from local version instead of bytes
	if commonCSV.file != nil {
		commonCSV.file.Close()
		log.Infof("file closed successfully")
	}

	commonCSV.eof = true
}

func (commonCSV *CommonCSV) IsEOF() bool {
	if commonCSV == nil {
		return true
	}

	return commonCSV.eof
}

// ParseHeaders parses and retrieves CSV headers
func (commonCSV *CommonCSV) ParseHeaders() (headers Headers, err error) {
	if commonCSV == nil {
		return
	}

	if commonCSV.Headers != nil {
		return commonCSV.Headers, nil
	}

	headers, err = commonCSV.Reader.Read()
	if err != nil {
		commonCSV.SetEOF()
		return nil, fmt.Errorf("error reading CSV headers: %w", err)
	}

	headers = Sanitize(headers, commonCSV.headerSanitizers)
	commonCSV.Headers = headers
	return headers, nil
}

// GetHeaders gets headers without re-parsing if already parsed
func (commonCSV *CommonCSV) GetHeaders() (Headers, error) {
	if commonCSV.Headers == nil {
		return commonCSV.ParseHeaders()
	}
	return commonCSV.Headers, nil
}

func (commonCSV *CommonCSV) ParseNextBatch(data interface{}) (err error) {
	if commonCSV == nil {
		err = fmt.Errorf("nil pointer reference")
		return
	}

	headers, err := commonCSV.ParseHeaders()
	if err != nil {
		return
	}

	records, err := commonCSV.ReadNextBatch()
	if err != nil {
		return
	}

	err = records.Unmarshal(headers, data)
	if err != nil {
		return
	}

	return
}

func (commonCSV *CommonCSV) S3Download(ctx context.Context, objectKey string, bucket string) ([]byte, error) {
	if objectKey == "" || bucket == "" {
		return nil, fmt.Errorf("objectkey or bucket is missing")
	}
	// Create an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.GetString(ctx, "s3.region")),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating AWS session: %v", err)
	}

	// Create an S3 service client
	svc := s3.New(sess)

	// Create an S3.GetObjectInput to specify the bucket and object to download
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	}

	// Use the S3 service client to download the object
	result, err := svc.GetObjectWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error downloading object from S3: %w", err)
	}

	defer result.Body.Close()

	objectData, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading object data: %w", err)
	}

	commonCSV.rawData = objectData
	return objectData, nil
}

func (commonCSV *CommonCSV) NewCSVReaderClient(ctx context.Context) (csvReader *csv.Reader) {
	if commonCSV == nil {
		return
	}

	// if file already downloaded use that to initialize reader
	if commonCSV.rawData != nil {
		reader := bytes.NewReader(commonCSV.rawData)
		csvReader = csv.NewReader(reader)
	}

	return
}

func (commonCSV *CommonCSV) ProcessSheet() (records Records, count int, err error) {
	if commonCSV == nil {
		return
	}

	if commonCSV.Headers == nil {
		commonCSV.ParseHeaders()
	}

	count = 0
	for {
		record, err := commonCSV.Reader.Read()
		if err == io.EOF {
			commonCSV.SetEOF()
			break
		} else if err != nil {
			return records, count, fmt.Errorf("error reading CSV record: %w", err)
		}

		count++
		records = append(records, record)
	}

	return records, count, nil
}
