package csv

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"os"
)

func (commonCSV *CommonCSV) InitializeS3CSVReader(ctx context.Context) (err error) {
	if commonCSV == nil {
		err = fmt.Errorf("nil pointer reference")
		return
	}

	// if file already downloaded use that to initialize reader
	if commonCSV.rawData != nil {
		reader := bytes.NewReader(commonCSV.rawData)
		commonCSV.Reader = csv.NewReader(reader)
		commonCSV.Reader.Comma = commonCSV.delimiter
		commonCSV.Reader.LazyQuotes = true
	}

	// Download S3 File and create csvReader client
	rawData, err := commonCSV.S3Download(ctx, commonCSV.fileInfo.ObjectKey, commonCSV.fileInfo.Bucket)
	if err != nil {
		return
	}

	reader := bytes.NewReader(rawData)
	commonCSV.Reader = csv.NewReader(reader)
	commonCSV.Reader.Comma = commonCSV.delimiter
	commonCSV.Reader.LazyQuotes = true

	return
}
func (commonCSV *CommonCSV) InitializeRawURLReader(ctx context.Context) (err error) {
	if commonCSV == nil {
		err = fmt.Errorf("nil pointer reference")
		return
	}

	// if file already downloaded use that to initialize reader
	if commonCSV.rawData != nil {
		reader := bytes.NewReader(commonCSV.rawData)
		commonCSV.Reader = csv.NewReader(reader)
		commonCSV.Reader.Comma = commonCSV.delimiter
		commonCSV.Reader.LazyQuotes = true
	}

	commonCSV.rawURLInfo.downloadedFilePath, err = downloadFile("", commonCSV.rawURLInfo.URL)
	if err != nil {
		return err
	}

	decoder, err := commonCSV.rawURLInfo.getDecoder()
	if err != nil {
		return fmt.Errorf("no valid decoder found")
	}

	commonCSV.rawURLInfo.decodedFilePath, err = decoder.decode()
	if err != nil {
		return err
	}

	commonCSV.localFilePath = commonCSV.rawURLInfo.decodedFilePath
	return commonCSV.InitializeLocalFileReader(ctx)

}

func (commonCSV *CommonCSV) InitializerRawBytes() error {

	reader := bytes.NewReader(commonCSV.rawData)
	commonCSV.Reader = csv.NewReader(reader)
	commonCSV.Reader.Comma = commonCSV.delimiter
	commonCSV.Reader.LazyQuotes = true
	return nil
}

func (commonCSV *CommonCSV) InitializeLocalFileReader(ctx context.Context) (err error) {

	// Open the CSV file
	file, err := os.Open(commonCSV.localFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Create a new CSV reader
	commonCSV.file = file // to be closed after reading
	commonCSV.Reader = csv.NewReader(file)
	commonCSV.Reader.Comma = commonCSV.delimiter
	commonCSV.Reader.LazyQuotes = true
	return
}

func (commonCSV *CommonCSV) InitializeReader(ctx context.Context) (err error) {
	source := commonCSV.GetSource()

	switch source {
	case S3:
		return commonCSV.InitializeS3CSVReader(ctx)
	case Local:
		return commonCSV.InitializeLocalFileReader(ctx)
	case RawBytes:
		return commonCSV.InitializerRawBytes()
	case RawURL:
		return commonCSV.InitializeRawURLReader(ctx)
	default:
		return commonCSV.InitializeS3CSVReader(ctx)
	}
}

func (commonCSV *CommonCSV) Close(ctx context.Context) (err error) {
	// If a file was opened (i.e. for Local or RawURL sources), close it.
	if commonCSV.file != nil {
		if cerr := commonCSV.file.Close(); cerr != nil {
			err = cerr
		}
		commonCSV.file = nil // clear reference after closing
	}

	// Clean up any temporary/local files for different sources.
	switch commonCSV.GetSource() {
	case Local, S3, S3Large, RawBytes:
		// No further cleanup required beyond closing the file.
		return
	case RawURL:
		// Remove temporary files created during raw URL processing
		if removeErr := os.Remove(commonCSV.rawURLInfo.decodedFilePath); removeErr != nil && err == nil {
			err = removeErr
		}

		if commonCSV.rawURLInfo.downloadedFilePath != commonCSV.rawURLInfo.decodedFilePath {
			if removeErr := os.Remove(commonCSV.rawURLInfo.downloadedFilePath); removeErr != nil && err == nil {
				err = removeErr
			}
		}
		return err
	default:
		return err
	}
}
