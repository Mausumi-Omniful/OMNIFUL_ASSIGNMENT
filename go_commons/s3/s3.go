package s3

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
	config2 "github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/log"
)

type AwsLogger struct {
}

type FileConfig struct {
	Path   string
	Bucket string
}

func (l *AwsLogger) Logf(classification logging.Classification, format string, args ...interface{}) {
	// todo: improve the logger and move out of this logrus
	switch classification {
	case logging.Warn:
		log.Warnf(format, args...)
	case logging.Debug:
		log.Debugf(format, args...)
	}
}

func NewDefaultAWSS3Client() (*awsS3.Client, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithLogger(&AwsLogger{}), // todo: allow to inject it as an option
	}

	if endpoint, ok := os.LookupEnv("AWS_S3_ENDPOINT"); ok {
		opts = append(opts, config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service string, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == s3.ServiceID {
				return aws.Endpoint{
					PartitionID:       "aws",
					URL:               endpoint,
					SigningRegion:     region,
					HostnameImmutable: true,
				}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})))
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		opts...,
	)
	if err != nil {
		return nil, err
	}

	// nrawssdk.InstrumentHandlers(&cfg.Handlers)
	// todo: use newrelic aws sdk v2.
	// todo: blocked on https://github.com/newrelic/go-agent/issues/288 https://github.com/newrelic/go-agent/issues/250
	client := awsS3.NewFromConfig(cfg)
	return client, nil
}

func GetURL(ctx context.Context, fileConfig FileConfig) string {
	return fmt.Sprintf("https://%s.%s/%s", fileConfig.Bucket, getBucketURL(ctx), fileConfig.Path)
}

func getBucketURL(ctx context.Context) string {
	if url, ok := os.LookupEnv("LOCAL_S3_BUCKET_URL"); ok {
		return url
	}

	region := config2.GetString(ctx, "s3.region")
	if region == "" {
		log.Errorf("s3 region is not set in config")
		return "s3.amazonaws.com"
	}

	return fmt.Sprintf("s3.%s.amazonaws.com", region)
}
