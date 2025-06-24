package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	commons3 "github.com/omniful/go_commons/s3"
)

type S3DownloaderImpl struct {
	client *s3.Client
}




// s3downloader
func NewS3Downloader(endpoint, region string) (*S3DownloaderImpl, error) {
	if endpoint != "" {
		os.Setenv("AWS_S3_ENDPOINT", endpoint)
	}
    client, err := commons3.NewDefaultAWSS3Client()
	if err != nil {
		return nil, fmt.Errorf("AWS client error: %w", err)
	}
    fmt.Println("S3 downloader ready")
	return &S3DownloaderImpl{client: client}, nil
}






// downloadfile
func (s *S3DownloaderImpl) DownloadFile(ctx context.Context, s3Path string) ([]byte, error) {
	bucket, key, err := s.parseS3Path(s3Path)
	if err != nil {
		return nil, err
	}

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("S3 download error: %w", err)
	}
	defer result.Body.Close()

	content, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	fmt.Printf("Downloaded %d bytes\n", len(content))
	return content, nil
}






// parses3path
func (s *S3DownloaderImpl) parseS3Path(s3Path string) (string, string, error) {
	if !strings.HasPrefix(s3Path, "s3://") {
		return "", "", fmt.Errorf("invalid path: %s", s3Path)
	}
	parts := strings.SplitN(strings.TrimPrefix(s3Path, "s3://"), "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid S3 path: %s", s3Path)
	}
	return parts[0], parts[1], nil
}
