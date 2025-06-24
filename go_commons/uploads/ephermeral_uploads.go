package uploads

import (
	"context"
	"fmt"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/s3"
	"time"
)

const (
	MetadataFilenameKey = "filename"
)

type TempUploadService struct {
	bucketName    string
	awsRegion     string
	s3Client      *awss3.Client
	presignClient *awss3.PresignClient
}

type TempURLRequest struct {
	Tenant  string `json:"tenant"`
	UseCase string `json:"usecase"`
}

type TempURLResponse struct {
	URL      string `json:"temp_url"`
	UploadID string `json:"upload_id"`
}

func NewTempUploadService(bucketName, awsRegion string) (*TempUploadService, error) {
	// Initialize S3 Client
	s3Client, err := s3.NewDefaultAWSS3Client()
	if err != nil {
		return nil, err
	}

	// Initialize Presign  s3 Client
	presignClient := awss3.NewPresignClient(
		s3Client,
		awss3.WithPresignExpires(24*time.Hour),
		awss3.WithPresignClientFromClientOptions(func() func(*awss3.Options) {
			return func(options *awss3.Options) {
				options.Region = awsRegion
			}
		}()),
	)

	return &TempUploadService{
		bucketName:    bucketName,
		awsRegion:     awsRegion,
		s3Client:      s3Client,
		presignClient: presignClient,
	}, nil
}

func (service *TempUploadService) GenerateTempURL(tenant, useCase, contentType, filename string) (*TempURLResponse, error) {
	// Load UTC timezone
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}

	// Generate unique upload ID
	currentTime := time.Now().In(loc)
	uploadID := fmt.Sprintf("%s/%s/%s/%s", tenant, useCase, currentTime.Format("2006-01-02/15/04/05"), uuid.New().String())

	// Define metadata for the S3 object
	metadata := map[string]string{
		MetadataFilenameKey: filename,
	}

	// Create input for presigned URL
	req := &awss3.PutObjectInput{
		Bucket:      aws.String(service.bucketName),
		Key:         aws.String(uploadID),
		ContentType: aws.String(contentType),
		Metadata:    metadata,
	}

	// Generate presigned URL
	resp, err := service.presignClient.PresignPutObject(context.TODO(), req)
	if err != nil {
		return nil, err
	}

	// Return response with presigned URL and upload ID
	return &TempURLResponse{
		URL:      resp.URL,
		UploadID: uploadID,
	}, nil
}
