# Package file_utilities

## Overview
The `file_utilities` package provides utilities for file operations with a specific focus on PDF file handling. It offers functionality for merging PDF files through AWS Lambda integration, making it easy to combine multiple PDFs into a single document.

## Features
- PDF Merging: Combine multiple PDFs into a single document
- AWS Lambda Integration: Utilizes AWS Lambda for processing PDF operations
- Error Handling: Comprehensive error reporting and failed PDF tracking
- URL-based Processing: Works with PDF files accessible via URLs

## Installation
```go
go get github.com/omniful/go_commons/file_utilities
```

## Components

### Requests
- `PdfMergeRequest`: Structure for PDF merge operations
  - `PdfURLs`: Array of PDFs to merge
  - `MergedPDFFilename`: Desired name for the merged PDF

### Responses
- `PdfMergeResponse`: Result of PDF merge operation
  - `FailedPdfs`: List of PDFs that failed during merging
  - `MergedPdfURL`: URL to the successfully merged PDF

## Usage Examples

### Merging PDFs
```go
package main

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/file_utilities"
	"github.com/omniful/go_commons/file_utilities/requests"
	"github.com/omniful/go_commons/lambda"
)

func main() {
	ctx := context.Background()
	
	// Initialize Lambda client
	lambdaClient := lambda.NewClient("your-aws-region")
	
	// Create merge request
	mergeRequest := requests.PdfMergeRequest{
		PdfURLs: requests.PdfURLs{
			{
				ID:  "doc1",
				URL: "https://example.com/doc1.pdf",
			},
			{
				ID:  "doc2",
				URL: "https://example.com/doc2.pdf",
			},
		},
		MergedPDFFilename: "merged_document.pdf",
	}
	
	// Execute merge operation
	response, err := file_utilities.PdfMerger(ctx, lambdaClient, mergeRequest)
	if err != nil {
		fmt.Printf("Error merging PDFs: %v\n", err)
		return
	}
	
	// Handle successful merge
	fmt.Printf("Merged PDF available at: %s\n", response.MergedPdfURL)
	
	// Check for any failed PDFs
	if len(response.FailedPdfs) > 0 {
		fmt.Println("Some PDFs failed to merge:")
		for _, failed := range response.FailedPdfs {
			fmt.Printf("- ID: %s, URL: %s\n", failed.ID, failed.URL)
		}
	}
}
```

## Error Handling
The package uses custom error types for different scenarios:
- `RequestNotValid`: When the lambda client is nil or no PDF URLs are provided
- `BadRequestError`: When the Lambda function execution fails

## Dependencies
- `github.com/omniful/go_commons/env`
- `github.com/omniful/go_commons/error`
- `github.com/omniful/go_commons/lambda`
- `github.com/omniful/go_commons/log`

## Notes
- Ensure AWS Lambda credentials are properly configured
- All PDF URLs must be publicly accessible
- The Lambda function "pdf_merger" must be deployed and available in your AWS environment
- PDF merge operations are asynchronous and processed through AWS Lambda

## Best Practices
1. Always check for failed PDFs in the response
2. Implement proper error handling
3. Use meaningful IDs for tracking PDFs in the merge request
4. Ensure PDF URLs are accessible to the Lambda function

## Contributing
Contributions are welcome! Please ensure your changes are well-tested and documented.
