package utils

import (
	"bytes"
	"context"
	"fmt"
	
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3UploaderImpl struct {
	client *s3.S3
	bucket string
}





func NewS3Uploader(bucket, endpoint, region string) (*S3UploaderImpl, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
		DisableSSL:       aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
	})
	if err != nil {
		return nil, fmt.Errorf("AWS session error: %w", err)
	}

	s3Client := s3.New(sess)

	_, err = s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		fmt.Println("Creating bucket:", bucket)
		_, err = s3Client.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return nil, fmt.Errorf("bucket creation error: %w", err)
		}
		fmt.Println("Bucket created:", bucket)
	} else {
		fmt.Println("Bucket exists:", bucket)
	}

	return &S3UploaderImpl{
		client: s3Client,
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

	_, err := s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(fileContent),
		ContentType: aws.String("text/csv"),
	})
	if err != nil {
		return "", fmt.Errorf("upload error: %w", err)
	}

	s3Path := fmt.Sprintf("s3://%s/%s", s.bucket, key)
	fmt.Println("Upload complete:", s3Path)
	return s3Path, nil
}





