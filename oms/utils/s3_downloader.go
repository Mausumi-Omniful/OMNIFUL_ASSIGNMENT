package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/omniful/go_commons/log"
)

type S3DownloaderImpl struct {
	client *s3.S3
}

// NewS3Downloader creates a new S3 downloader with LocalStack configuration
func NewS3Downloader(endpoint, region string) (*S3DownloaderImpl, error) {
	// Create AWS session for LocalStack
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint), // localstack endpoint
		S3ForcePathStyle: aws.Bool(true),       // required for localstack
		DisableSSL:       aws.Bool(true),       // disable SSL for localstack
		Credentials: credentials.NewStaticCredentials(
			"test", "test", "", // dummy creds for localstack
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	log.Infof("S3 Downloader initialized for endpoint: %s, region: %s", endpoint, region)

	return &S3DownloaderImpl{
		client: s3.New(sess),
	}, nil
}

// DownloadFile downloads a file from S3 and returns its content as bytes
func (s *S3DownloaderImpl) DownloadFile(ctx context.Context, s3Path string) ([]byte, error) {
	// Parse S3 path: s3://bucket-name/key
	bucket, key, err := s.parseS3Path(s3Path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse S3 path: %w", err)
	}

	log.Infof("Downloading file from S3: bucket=%s, key=%s", bucket, key)

	// Download file from S3
	result, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file from S3: %w", err)
	}
	defer result.Body.Close()

	// Read the file content
	content, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	log.Infof("File downloaded successfully: size=%d bytes", len(content))

	// Log content preview for debugging
	if len(content) > 0 {
		preview := string(content)
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		log.Infof("File content preview: %s", preview)
	} else {
		log.Warn("Downloaded file is empty!")
	}

	return content, nil
}

// DownloadFileToPath downloads a file from S3 and saves it to a local path
func (s *S3DownloaderImpl) DownloadFileToPath(ctx context.Context, s3Path, localPath string) error {
	// Download file content
	content, err := s.DownloadFile(ctx, s3Path)
	if err != nil {
		return err
	}

	// Create local file
	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer file.Close()

	// Write content to local file
	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write to local file: %w", err)
	}

	log.Infof("File saved to local path: %s", localPath)
	return nil
}

// FileExists checks if a file exists in S3
func (s *S3DownloaderImpl) FileExists(ctx context.Context, s3Path string) (bool, error) {
	bucket, key, err := s.parseS3Path(s3Path)
	if err != nil {
		return false, fmt.Errorf("failed to parse S3 path: %w", err)
	}

	_, err = s.client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		// Check if it's a "not found" error
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "NoSuchKey") {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	return true, nil
}

// GetFileSize returns the size of a file in S3
func (s *S3DownloaderImpl) GetFileSize(ctx context.Context, s3Path string) (int64, error) {
	bucket, key, err := s.parseS3Path(s3Path)
	if err != nil {
		return 0, fmt.Errorf("failed to parse S3 path: %w", err)
	}

	result, err := s.client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get file size: %w", err)
	}

	if result.ContentLength == nil {
		return 0, fmt.Errorf("content length is nil")
	}

	return *result.ContentLength, nil
}

// parseS3Path parses an S3 path and returns bucket and key
func (s *S3DownloaderImpl) parseS3Path(s3Path string) (string, string, error) {
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

// ValidateS3Path validates if the S3 path format is correct
func (s *S3DownloaderImpl) ValidateS3Path(s3Path string) error {
	_, _, err := s.parseS3Path(s3Path)
	return err
}
