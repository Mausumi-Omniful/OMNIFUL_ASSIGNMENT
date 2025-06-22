package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/omniful/go_commons/csv"
	"github.com/omniful/go_commons/log"
)

// CSVRow represents a single row from the CSV file
type CSVRow struct {
	SKU       string `json:"sku"`
	Location  string `json:"location"`
	TenantID  string `json:"tenant_id"`
	SellerID  string `json:"seller_id"`
	RowNumber int    `json:"row_number"`
}

// CSVParseResult represents the result of CSV parsing
type CSVParseResult struct {
	TotalRows    int      `json:"total_rows"`
	ValidRows    int      `json:"valid_rows"`
	InvalidRows  int      `json:"invalid_rows"`
	ValidData    []CSVRow `json:"valid_data"`
	InvalidData  []CSVRow `json:"invalid_data"`
	Headers      []string `json:"headers"`
	ErrorRows    []int    `json:"error_rows"`
	ErrorMessage string   `json:"error_message,omitempty"`
}

// CSVParser handles CSV file parsing using go_commons
type CSVParser struct {
	batchSize int
}

// NewCSVParser creates a new CSV parser instance
func NewCSVParser(batchSize int) *CSVParser {
	if batchSize <= 0 {
		batchSize = 100 // default batch size
	}
	return &CSVParser{
		batchSize: batchSize,
	}
}

// ParseCSVFromBytes parses CSV content from bytes using go_commons
func (p *CSVParser) ParseCSVFromBytes(ctx context.Context, csvData []byte) (*CSVParseResult, error) {
	log.Infof("Starting CSV parsing with batch size: %d", p.batchSize)

	result := &CSVParseResult{
		ValidData:   make([]CSVRow, 0),
		InvalidData: make([]CSVRow, 0),
		ErrorRows:   make([]int, 0),
	}

	// Create CSV reader with raw bytes
	csvReader, err := csv.NewCommonCSV(
		csv.WithRawData(csvData),
		csv.WithSource(csv.RawBytes),
		csv.WithBatchSize(p.batchSize),
		csv.WithHeaderSanitizers(csv.SanitizeSpace),
		csv.WithDataRowSanitizers(csv.SanitizeSpace),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create CSV reader: %w", err)
	}

	// Initialize the reader
	err = csvReader.InitializeReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CSV reader: %w", err)
	}

	// Parse headers
	headers, err := csvReader.ParseHeaders()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV headers: %w", err)
	}

	result.Headers = headers
	log.Infof("CSV Headers: %v", headers)

	// Validate headers
	if err := p.validateHeaders(headers); err != nil {
		return nil, fmt.Errorf("invalid CSV headers: %w", err)
	}

	// Process CSV in batches
	rowNumber := 1 // Start from 1 since we already read headers
	for !csvReader.IsEOF() {
		var batchData []CSVRow
		err := csvReader.ParseNextBatch(&batchData)
		if err != nil {
			log.WithError(err).Errorf("Failed to parse batch starting at row %d", rowNumber+1)
			result.ErrorMessage = fmt.Sprintf("Failed to parse batch at row %d: %v", rowNumber+1, err)
			break
		}

		// Process each row in the batch
		for i, row := range batchData {
			row.RowNumber = rowNumber + i + 1

			// Validate the row
			if err := p.validateRow(row); err != nil {
				log.WithError(err).Warnf("Invalid row %d: %v", row.RowNumber, err)
				result.InvalidData = append(result.InvalidData, row)
				result.InvalidRows++
				result.ErrorRows = append(result.ErrorRows, row.RowNumber)
			} else {
				result.ValidData = append(result.ValidData, row)
				result.ValidRows++
			}
			result.TotalRows++
		}

		rowNumber += len(batchData)
		log.Infof("Processed batch: %d rows, Total processed: %d", len(batchData), result.TotalRows)
	}

	// Close the reader
	err = csvReader.Close(ctx)
	if err != nil {
		log.WithError(err).Warn("Failed to close CSV reader")
	}

	log.Infof("CSV parsing completed - Total: %d, Valid: %d, Invalid: %d",
		result.TotalRows, result.ValidRows, result.InvalidRows)

	return result, nil
}

// ParseCSVFromS3Path parses CSV file directly from S3 using go_commons
func (p *CSVParser) ParseCSVFromS3Path(ctx context.Context, s3Path string) (*CSVParseResult, error) {
	// Parse S3 path to get bucket and key
	bucket, key, err := p.parseS3Path(s3Path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse S3 path: %w", err)
	}

	log.Infof("Starting CSV parsing from S3: bucket=%s, key=%s", bucket, key)

	// Create CSV reader with S3 source
	csvReader, err := csv.NewCommonCSV(
		csv.WithFileInfo(key, bucket),
		csv.WithSource(csv.S3),
		csv.WithBatchSize(p.batchSize),
		csv.WithHeaderSanitizers(csv.SanitizeSpace),
		csv.WithDataRowSanitizers(csv.SanitizeSpace),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create CSV reader: %w", err)
	}

	// Initialize the reader (this will download from S3)
	err = csvReader.InitializeReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CSV reader: %w", err)
	}

	// Parse headers
	headers, err := csvReader.ParseHeaders()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV headers: %w", err)
	}

	result := &CSVParseResult{
		Headers:     headers,
		ValidData:   make([]CSVRow, 0),
		InvalidData: make([]CSVRow, 0),
		ErrorRows:   make([]int, 0),
	}

	log.Infof("CSV Headers: %v", headers)

	// Validate headers
	if err := p.validateHeaders(headers); err != nil {
		return nil, fmt.Errorf("invalid CSV headers: %w", err)
	}

	// Process CSV in batches
	rowNumber := 1 // Start from 1 since we already read headers
	for !csvReader.IsEOF() {
		var batchData []CSVRow
		err := csvReader.ParseNextBatch(&batchData)
		if err != nil {
			log.WithError(err).Errorf("Failed to parse batch starting at row %d", rowNumber+1)
			result.ErrorMessage = fmt.Sprintf("Failed to parse batch at row %d: %v", rowNumber+1, err)
			break
		}

		// Process each row in the batch
		for i, row := range batchData {
			row.RowNumber = rowNumber + i + 1

			// Validate the row
			if err := p.validateRow(row); err != nil {
				log.WithError(err).Warnf("Invalid row %d: %v", row.RowNumber, err)
				result.InvalidData = append(result.InvalidData, row)
				result.InvalidRows++
				result.ErrorRows = append(result.ErrorRows, row.RowNumber)
			} else {
				result.ValidData = append(result.ValidData, row)
				result.ValidRows++
			}
			result.TotalRows++
		}

		rowNumber += len(batchData)
		log.Infof("Processed batch: %d rows, Total processed: %d", len(batchData), result.TotalRows)
	}

	// Close the reader
	err = csvReader.Close(ctx)
	if err != nil {
		log.WithError(err).Warn("Failed to close CSV reader")
	}

	log.Infof("CSV parsing completed - Total: %d, Valid: %d, Invalid: %d",
		result.TotalRows, result.ValidRows, result.InvalidRows)

	return result, nil
}

// validateHeaders validates that the CSV has the required headers
func (p *CSVParser) validateHeaders(headers []string) error {
	requiredHeaders := []string{"sku", "location", "tenant_id", "seller_id"}
	headerMap := make(map[string]bool)

	// Convert headers to lowercase for case-insensitive comparison
	for _, header := range headers {
		headerMap[strings.ToLower(strings.TrimSpace(header))] = true
	}

	// Check for required headers
	for _, required := range requiredHeaders {
		if !headerMap[required] {
			return fmt.Errorf("missing required header: %s", required)
		}
	}

	return nil
}

// validateRow validates a single CSV row
func (p *CSVParser) validateRow(row CSVRow) error {
	// Check for empty required fields
	if strings.TrimSpace(row.SKU) == "" {
		return fmt.Errorf("SKU is empty")
	}
	if strings.TrimSpace(row.Location) == "" {
		return fmt.Errorf("location is empty")
	}
	if strings.TrimSpace(row.TenantID) == "" {
		return fmt.Errorf("tenant_id is empty")
	}
	if strings.TrimSpace(row.SellerID) == "" {
		return fmt.Errorf("seller_id is empty")
	}

	// Add more validation rules as needed
	// For example: SKU format, location validation, etc.

	return nil
}

// parseS3Path parses an S3 path and returns bucket and key
func (p *CSVParser) parseS3Path(s3Path string) (string, string, error) {
	// Remove s3:// prefix
	if !strings.HasPrefix(s3Path, "s3://") {
		return "", "", fmt.Errorf("invalid S3 path format, must start with 's3://': %s", s3Path)
	}

	// Remove s3:// prefix
	path := strings.TrimPrefix(s3Path, "s3://")

	// Split by first slash to separate bucket and key
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid S3 path format, must be 's3://bucket/key': %s", s3Path)
	}

	bucket := parts[0]
	key := parts[1]

	if bucket == "" {
		return "", "", fmt.Errorf("bucket name cannot be empty: %s", s3Path)
	}

	if key == "" {
		return "", "", fmt.Errorf("key cannot be empty: %s", s3Path)
	}

	return bucket, key, nil
}
