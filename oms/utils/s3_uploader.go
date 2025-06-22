package utils

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/omniful/go_commons/log"
)

type S3UploaderImpl struct {
	client *s3.S3
	bucket string
}

func NewS3Uploader(bucket, endpoint, region string) (*S3UploaderImpl, error) {
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

	s3Client := s3.New(sess)

	// Check if bucket exists and create it if it doesn't
	_, err = s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		// Bucket doesn't exist, create it
		log.Infof("Creating S3 bucket: %s", bucket)

		_, err = s3Client.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})

		if err != nil {
			return nil, fmt.Errorf("failed to create S3 bucket: %w", err)
		}

		log.Infof("S3 bucket created: %s", bucket)
	} else {
		log.Infof("S3 bucket already exists: %s", bucket)
	}

	return &S3UploaderImpl{
		client: s3Client,
		bucket: bucket,
	}, nil
}

func (s *S3UploaderImpl) UploadFile(ctx context.Context, fileContent []byte, filename string) (string, error) {
	// Generate unique key with timestamp
	timestamp := time.Now().Unix()
	key := fmt.Sprintf("csv-uploads/%d-%s", timestamp, filename)

	log.Infof("Uploading file to S3: bucket=%s, key=%s, size=%d bytes", s.bucket, key, len(fileContent))

	// Debug: Log content preview before upload
	if len(fileContent) > 0 {
		preview := string(fileContent)
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		log.Infof("S3 upload content preview: %s", preview)
	} else {
		log.Warn("S3 upload: fileContent is empty!")
	}

	// Upload file to S3
	_, err := s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(fileContent),
		ContentType: aws.String("text/csv"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Generate S3 path
	s3Path := fmt.Sprintf("s3://%s/%s", s.bucket, key)
	log.Infof("File uploaded successfully: %s", s3Path)

	return s3Path, nil
}

// GeneratePresignedURL generates a pre-signed URL for downloading files from S3
func (s *S3UploaderImpl) GeneratePresignedURL(ctx context.Context, s3Path string, expirySeconds int64) (string, error) {
	// Parse S3 path to extract bucket and key
	bucket, key, err := s.parseS3Path(s3Path)
	if err != nil {
		return "", fmt.Errorf("failed to parse S3 path: %w", err)
	}

	log.Infof("Generating pre-signed URL for S3: bucket=%s, key=%s, expiry=%d seconds", bucket, key, expirySeconds)

	// Create pre-signed URL request
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	// Generate pre-signed URL
	url, err := req.Presign(time.Duration(expirySeconds) * time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to generate pre-signed URL: %w", err)
	}

	log.Infof("Pre-signed URL generated successfully")
	return url, nil
}

// parseS3Path parses an S3 path to extract bucket and key
func (s *S3UploaderImpl) parseS3Path(s3Path string) (string, string, error) {
	// Expected format: s3://bucket/key
	if len(s3Path) < 6 || s3Path[:5] != "s3://" {
		return "", "", fmt.Errorf("invalid S3 path format: %s", s3Path)
	}

	// Remove s3:// prefix
	path := s3Path[5:]

	// Split by first slash to separate bucket and key
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid S3 path format: %s", s3Path)
	}

	return parts[0], parts[1], nil
}
