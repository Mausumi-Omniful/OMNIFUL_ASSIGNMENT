package utils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	commons3 "github.com/omniful/go_commons/s3"
)

type S3UploaderImpl struct {
	client *s3.Client
	bucket string
}

func NewS3Uploader(bucket, endpoint, region string) (*S3UploaderImpl, error) {
	if endpoint != "" {
		os.Setenv("AWS_S3_ENDPOINT", endpoint)
	}

	if os.Getenv("AWS_REGION") == "" && region != "" {
		os.Setenv("AWS_REGION", region)
	}
	if os.Getenv("AWS_REGION") == "" {
		return nil, fmt.Errorf("AWS region must be set")
	}

	client, err := commons3.NewDefaultAWSS3Client()
	if err != nil {
		return nil, fmt.Errorf("AWS client error: %w", err)
	}

	_, err = client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: &bucket,
	})
	if err != nil {
		fmt.Println("Creating bucket:", bucket)
		_, err = client.CreateBucket(context.Background(), &s3.CreateBucketInput{
			Bucket: &bucket,
		})
		if err != nil {
			return nil, fmt.Errorf("bucket creation error: %w", err)
		}
		fmt.Println("Bucket created:",bucket)
	} else {
		fmt.Println("Bucket exists:",bucket)
	}

	return &S3UploaderImpl{
		client: client,
		bucket: bucket,
	}, nil
}





func (s *S3UploaderImpl) UploadFile(ctx context.Context, fileContent []byte, filename string) (string, error) {
	timestamp := time.Now().Unix()
	key := fmt.Sprintf("csv-uploads/%d-%s", timestamp, filename)

	fmt.Printf("Uploading to S3: %s/%s (%d bytes)\n", s.bucket, key, len(fileContent))

	if len(fileContent) > 0 {
		preview := string(fileContent)
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		fmt.Println("Preview:", preview)
	} else {
		fmt.Println("Warning: file content is empty")
	}

	contentType := "text/csv"
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &key,
		Body:        bytes.NewReader(fileContent),
		ContentType: &contentType,
	})
	if err != nil {
		return "", fmt.Errorf("upload error: %w", err)
	}

	s3Path := fmt.Sprintf("s3://%s/%s", s.bucket, key)
	fmt.Println("Upload complete:", s3Path)
	return s3Path, nil
}