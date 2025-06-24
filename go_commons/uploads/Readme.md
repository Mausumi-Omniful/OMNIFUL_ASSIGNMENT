# Package uploads

## Overview
The uploads package provides a robust solution for handling ephemeral file uploads and downloads using AWS S3. It implements a two-stage file handling system:

1. **Stage 1 - Temporary Storage**: Files are initially uploaded to a single public bucket using pre-signed URLs
2. **Stage 2 - Processing/Permanent Storage**: Files can then be either:
   - Downloaded for immediate processing (e.g., bulk order processing)
   - Transferred to secure private buckets for permanent storage (e.g., SKU images, invoices)

This approach centralizes file uploads while allowing different services to handle files according to their specific needs.

## Features
- Pre-signed URL generation for direct browser-to-S3 uploads
- Temporary file storage in a centralized public bucket
- File download capabilities for immediate processing
- Secure file transfer to private buckets for permanent storage
- Metadata handling and ETag management
- Support for tenant and use-case based file organization

## Installation
```go
go get github.com/omniful/go_commons/uploads
```

## Components

### TempUploadService
Generates pre-signed URLs that allow frontend applications to upload files directly to S3, bypassing your application server.

```go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/uploads"
)

func main() {
	// Initialize the upload service with the public temporary bucket
	service, err := uploads.NewTempUploadService("public-temp-bucket", "us-west-2")
	if err != nil {
		panic(err)
	}

	// Generate a temporary upload URL for frontend use
	response, err := service.GenerateTempURL(
		"marketplace",        // tenant
		"bulk-orders",       // useCase
		"text/csv",          // contentType
		"orders.csv",        // filename
	)
	if err != nil {
		panic(err)
	}

	// The response contains:
	// - URL: Pre-signed URL for frontend to upload directly to S3
	// - UploadID: Unique identifier to reference this file later
	fmt.Printf("Upload URL: %s\n", response.URL)
	fmt.Printf("Upload ID: %s\n", response.UploadID)
}
```

### EphemeralClient
Handles post-upload operations: downloading for processing or transferring to permanent storage.

```go
package main

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/uploads"
	"os"
)

func main() {
	// Initialize client with the public temporary bucket
	client, err := uploads.NewEphemeralDownloadClient("public-temp-bucket")
	if err != nil {
		panic(err)
	}

	// Example 1: Download file for immediate processing (e.g., bulk order CSV)
	// The uploadID is the one received from TempUploadService.GenerateTempURL
	uploadID := "marketplace/bulk-orders/2024-02-06/15/04/05/uuid123"
	
	file, err := os.Create("orders.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = client.DownloadObject(context.TODO(), &uploads.DownloadObjectInput{
		EphemeralUploadID: uploadID,  // Required: Upload ID from pre-signed URL generation
		File:             file,
	})
	if err != nil {
		panic(err)
	}

	// Example 2: Transfer file to permanent storage (e.g., SKU images)
	response, err := client.Upload(
		context.TODO(),
		"private-permanent-bucket",    // Target bucket for permanent storage
		uploadID,                      // Upload ID from pre-signed URL generation
		"skus/images",                 // Permanent storage directory
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Permanent file URL: %s\n", response.URL)
	fmt.Printf("Media ID: %s\n", response.MediaID)
}
```

## Common Use Cases

1. **Bulk Order Processing**
   - Frontend uploads CSV to temporary bucket using pre-signed URL
   - Backend downloads file immediately for processing
   - File can be discarded after processing

2. **SKU Image Storage**
   - Frontend uploads image to temporary bucket
   - Backend transfers file to permanent private bucket
   - Image remains accessible indefinitely

3. **Order Invoices**
   - Invoice PDF uploaded to temporary bucket
   - Transferred to secure permanent storage
   - Accessible for future reference

4. **GRN Attachments**
   - Documents uploaded to temporary storage
   - Moved to permanent storage for record-keeping
   - Secure access for authorized users

## Response Types

### TempURLResponse
```go
type TempURLResponse struct {
	URL      string `json:"temp_url"`  // Pre-signed URL for frontend upload
	UploadID string `json:"upload_id"` // ID to reference this file later
}
```

### UploadResponse
```go
type UploadResponse struct {
	MediaID       string // Unique permanent media identifier
	Md5Hash       string // MD5 hash for integrity verification
	FileExtension string // File extension
	Directory     string // Permanent storage location
	Bucket        string // Target bucket name
	URL           string // Permanent access URL
}
```

## Best Practices
1. Always store the `uploadID` from TempURLResponse - it's required for subsequent operations
2. Use appropriate content types when generating upload URLs
3. Implement proper error handling for S3 operations
4. Clean up temporary files after processing
5. Use secure private buckets for permanent storage
6. Consider implementing retry mechanisms for failed operations

## Environment Variables
- `LOCAL_S3_BUCKET_URL`: Custom S3 bucket URL (optional, defaults to s3.amazonaws.com)

## Notes
- The package requires proper AWS credentials configuration
- Pre-signed URLs are valid for 24 hours by default
- File paths in temporary storage are organized by tenant/use-case/timestamp
- The package automatically handles metadata and ETag management
- Temporary storage should be regularly cleaned up to prevent accumulation of unused files
