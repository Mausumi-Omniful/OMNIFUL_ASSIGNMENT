package uploads

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	awss3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/omniful/go_commons/s3"
	"github.com/spf13/cast"
)

type UploadResponse struct {
	MediaID       string
	Md5Hash       string
	FileExtension string
	Directory     string
	Bucket        string
	URL           string
}

type EphemeralClient struct {
	s3client     *awss3.Client
	sourceBucket string
}

type CopyObjectInput struct {
	EphemeralUploadID string
	TargetBucket      string
	TargetPath        string
	Region            string
}

type DownloadObjectInput struct {
	EphemeralUploadID string
	File              *os.File
}

func NewEphemeralDownloadClient(sourceBucket string) (*EphemeralClient, error) {
	s3Client, err := s3.NewDefaultAWSS3Client()
	if err != nil {
		return nil, err
	}

	return &EphemeralClient{
		s3client:     s3Client,
		sourceBucket: sourceBucket,
	}, nil
}

func (ef *EphemeralClient) DownloadObject(ctx context.Context, param *DownloadObjectInput) error {
	resp, err := ef.s3client.GetObject(
		ctx,
		&awss3.GetObjectInput{
			Bucket: aws.String(ef.sourceBucket),
			Key:    aws.String(param.EphemeralUploadID),
		})
	if err != nil {
		return err
	}

	s3objectBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	f := param.File
	_, err = f.Write(s3objectBytes)
	if err != nil {
		return err
	}

	return nil
}

func (ef *EphemeralClient) CopyObject(ctx context.Context, param *CopyObjectInput) (uploadedURL string, err error) {
	_, err = ef.s3client.CopyObject(
		ctx, &awss3.CopyObjectInput{
			Bucket:            aws.String(param.TargetBucket),
			CopySource:        aws.String(ef.sourceBucket + "/" + param.EphemeralUploadID),
			Key:               aws.String(param.TargetPath),
			MetadataDirective: awss3types.MetadataDirectiveCopy,
		})
	if err != nil {
		return
	}

	uploadedURL = s3.GetURL(ctx, s3.FileConfig{
		Bucket: param.TargetBucket,
		Path:   param.TargetPath,
	})

	return
}

func (ef *EphemeralClient) GetHeadObject(ctx context.Context, uploadID string) (*awss3.HeadObjectOutput, error) {
	headObject := &awss3.HeadObjectInput{
		Bucket: aws.String(ef.sourceBucket),
		Key:    aws.String(uploadID),
	}

	response, err := ef.s3client.HeadObject(ctx, headObject)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (ef *EphemeralClient) GetETagWithExtension(ctx context.Context, uploadID string) (string, string, error) {
	response, err := ef.GetHeadObject(ctx, uploadID)
	if err != nil {
		return "", "", fmt.Errorf("error in getting hash for the uploaded files error %v", err)
	}

	etag := strings.Trim(cast.ToString(response.ETag), "\"")
	extension := strings.Replace(cast.ToString(response.ContentType), "image/", "", -1)
	return etag, extension, nil
}

// Upload
// directory will be without slash example - skus/images
// uploadID will be the one that is coming from api gateway
// targetBucket will be the one where the upload need to be stored
func (ef *EphemeralClient) Upload(ctx context.Context, targetBucket, uploadID, directory string) (*UploadResponse, error) {
	md5hash, extension, err := ef.GetETagWithExtension(ctx, uploadID)
	if err != nil {
		return nil, err
	}

	mediaID := getMediaID(md5hash, extension)
	filePath := getFilePath(directory, mediaID)

	uploadedURL, copyErr := ef.CopyObject(ctx, &CopyObjectInput{
		EphemeralUploadID: uploadID,
		TargetBucket:      targetBucket,
		TargetPath:        filePath,
	})
	if copyErr != nil {
		return nil, err
	}

	result := &UploadResponse{
		MediaID:       mediaID,
		Md5Hash:       md5hash,
		FileExtension: extension,
		Directory:     directory,
		Bucket:        targetBucket,
		URL:           uploadedURL,
	}

	return result, nil
}

func getMediaID(hash, extension string) string {
	mediaID := fmt.Sprintf("%s%d.%s", hash, time.Now().Unix(), extension)
	return mediaID
}

func getFilePath(directory, fileName string) string {
	return directory + "/" + fileName
}
