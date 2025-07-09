package utils

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestS3UploaderImpl_UploadFile_Integration(t *testing.T) {
	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	os.Setenv("AWS_REGION", "us-east-1")

	bucket := "test-bucket"
	endpoint := "http://localhost:4566"
	region := "us-east-1"

	uploader, err := NewS3Uploader(bucket, endpoint, region)
	assert.NoError(t, err)
	assert.NotNil(t, uploader)

	ctx := context.Background()
	fileContent := []byte("integration test csv content")
	filename := "integration_test.csv"

	path, err := uploader.UploadFile(ctx, fileContent, filename)
	assert.NoError(t, err)
	assert.Contains(t, path, "s3://"+bucket+"/csv-uploads/")
	assert.Contains(t, path, filename)
}
