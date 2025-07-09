package utils

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQSPublisherImpl_PublishS3Path_Integration(t *testing.T) {
	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	os.Setenv("AWS_DEFAULT_REGION", "us-east-1")

	queueName := "test-queue"
	endpoint := "http://localhost:4566"
	region := "us-east-1"

	publisher, err := NewSQSPublisher(queueName, endpoint, region)
	assert.NoError(t, err)
	assert.NotNil(t, publisher)

	ctx := context.Background()
	s3Path := "s3://test-bucket/csv-uploads/integration_test.csv"

	err = publisher.PublishS3Path(ctx, s3Path)
	assert.NoError(t, err)
}
