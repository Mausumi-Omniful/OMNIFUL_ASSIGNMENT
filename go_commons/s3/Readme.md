# Package s3

## Overview
The s3 package provides an abstraction layer for interacting with Amazon S3. It simplifies operations like uploading, downloading, and managing buckets to facilitate cloud storage integration.

## Key Components
- S3 Client: Manages authentication and connection to S3.
- File Operations: Upload, download, and delete files.
- Bucket Management: Create, list, and delete buckets.

## Usage Example
~~~go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/s3"
)

func main() {
	client, err := s3.NewClient(s3.Config{
		// Configuration options
	})
	if err != nil {
		fmt.Println("Error creating S3 client:", err)
		return
	}
	err = client.Upload("bucket-name", "file.txt", []byte("file content"))
	if err != nil {
		fmt.Println("Upload failed:", err)
	} else {
		fmt.Println("File uploaded successfully!")
	}
}
~~~

## Notes
- Extendable for advanced S3 operations.
