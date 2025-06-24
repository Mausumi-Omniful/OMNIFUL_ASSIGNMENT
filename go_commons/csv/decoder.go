package csv

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Define the possible encodings
const (
	EncodingCSV  = "csv"
	EncodingGzip = "gzip"
	EncodingZip  = "zip"
)

type Decoder interface {
	decode() (string, error)
}

type GzipDecoder struct {
	filePath string
}

type ZipDecoder struct {
	filePath string
}

// CSVDecoder struct for default handling
type CSVDecoder struct {
	filePath string
}

// Function to get the appropriate decoder based on encoding
func (info *RawURLInfo) getDecoder() (decoder Decoder, err error) {
	if info == nil {
		return decoder, fmt.Errorf("nil pointer reference")
	}

	switch info.Encoding {
	case EncodingGzip:
		return &GzipDecoder{filePath: info.downloadedFilePath}, nil
	case EncodingZip:
		return &ZipDecoder{filePath: info.downloadedFilePath}, nil
	default:
		return &CSVDecoder{filePath: info.downloadedFilePath}, nil
	}
}

func (d *GzipDecoder) decode() (string, error) {
	// Open the gzip file
	file, err := os.Open(d.filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a gzip reader
	gz, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gz.Close()

	// Calculate the name of the uncompressed file
	baseName := strings.TrimSuffix(filepath.Base(d.filePath), ".gz")
	csvFileName := filepath.Join(filepath.Dir(d.filePath), baseName)

	// Create a new file to write the uncompressed content
	outFile, err := os.Create(csvFileName)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	// Write the file content in chunks
	buf := make([]byte, 32*1024) // 32KB chunks
	for {
		n, err := gz.Read(buf)
		if n > 0 {
			_, writeErr := outFile.Write(buf[:n])
			if writeErr != nil {
				return "", writeErr
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
	}

	return csvFileName, nil
}

func (d *ZipDecoder) decode() (string, error) {
	// Open the zip file
	r, err := zip.OpenReader(d.filePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	// Assuming the zip file contains only one file, which is the CSV
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()

		// Create a new file to write the content
		csvFileName := filepath.Join(filepath.Dir(d.filePath), f.Name)
		outFile, err := os.Create(csvFileName)
		if err != nil {
			return "", err
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return "", err
		}

		return csvFileName, nil
	}
	return "", fmt.Errorf("no files found in the zip archive")
}

// CSVDecoder decode function (default)
func (d *CSVDecoder) decode() (string, error) {
	// For CSV, just return the file path itself
	return d.filePath, nil
}
