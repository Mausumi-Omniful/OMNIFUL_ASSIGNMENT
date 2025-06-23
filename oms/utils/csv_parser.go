package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/omniful/go_commons/csv"
)




type CSVRow struct {
	SKU       string `json:"sku"`
	Location  string `json:"location"`
	TenantID  string `json:"tenant_id"`
	SellerID  string `json:"seller_id"`
	RowNumber int    `json:"row_number"`
}





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

type CSVParser struct {
	batchSize int
}





func NewCSVParser(batchSize int) *CSVParser {
	if batchSize <= 0 {
		batchSize = 100 
	}
	return &CSVParser{
		batchSize: batchSize,
	}
}




// parsecsvfrombytes
func (p *CSVParser) ParseCSVFromBytes(ctx context.Context, csvData []byte) (*CSVParseResult, error) {
	fmt.Printf("Starting CSV parsing with batch size: %d\n", p.batchSize)

	result := &CSVParseResult{
		ValidData:   make([]CSVRow, 0),
		InvalidData: make([]CSVRow, 0),
		ErrorRows:   make([]int, 0),
	}

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

	if err = csvReader.InitializeReader(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize CSV reader: %w", err)
	}

	headers, err := csvReader.ParseHeaders()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV headers: %w", err)
	}

	result.Headers = headers
	fmt.Printf("CSV Headers: %v\n", headers)

	if err := p.validateHeaders(headers); err != nil {
		return nil, fmt.Errorf("invalid CSV headers: %w", err)
	}

	rowNumber := 1
	for !csvReader.IsEOF() {
		var batchData []CSVRow
		err := csvReader.ParseNextBatch(&batchData)
		if err != nil {
			result.ErrorMessage = fmt.Sprintf("Failed to parse batch at row %d: %v", rowNumber+1, err)
			break
		}

		for i, row := range batchData {
			row.RowNumber = rowNumber + i + 1
			if err := p.validateRow(row); err != nil {
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
		fmt.Printf("Processed batch: %d rows, Total processed: %d\n", len(batchData), result.TotalRows)
	}

	_ = csvReader.Close(ctx)

	fmt.Printf("CSV parsing completed - Total: %d, Valid: %d, Invalid: %d\n",
		result.TotalRows, result.ValidRows, result.InvalidRows)

	return result, nil
}













// validateheaders
func (p *CSVParser) validateHeaders(headers []string) error {
	required := []string{"sku", "location", "tenant_id", "seller_id"}
	m := make(map[string]bool)
	for _, h := range headers {
		m[strings.ToLower(strings.TrimSpace(h))] = true
	}
	for _, r := range required {
		if !m[r] {
			return fmt.Errorf("missing required header: %s", r)
		}
	}
	return nil
}






// validaterow
func (p *CSVParser) validateRow(row CSVRow) error {
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
	return nil
}





