package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/env"
	"github.com/omniful/go_commons/log"
)

type CommonCSVWriter struct {
	Writer          *csv.Writer
	Headers         Headers
	file            *os.File
	destination     Destination
	isHeaderWritten bool
	shouldAddBOM    bool
	rows            int
}

func NewCommonCSVWriter(options ...WriterOptionFunc) (*CommonCSVWriter, error) {
	opts := &WriterOptions{}

	for _, option := range options {
		option(opts)
	}

	return &CommonCSVWriter{
		Writer:          opts.Writer,
		Headers:         opts.Headers,
		file:            opts.file,
		destination:     opts.destination,
		isHeaderWritten: opts.isHeaderWritten,
		shouldAddBOM:    opts.shouldAddBOM,
	}, nil
}

func (writer *CommonCSVWriter) IsHeaderWritten() bool {
	return writer.isHeaderWritten
}

func (writer *CommonCSVWriter) SetHeaderWritten() {
	writer.isHeaderWritten = true
}

func (writer *CommonCSVWriter) GetDestination() (destination Destination) {
	if writer == nil {
		return
	}

	return writer.destination
}

func (writer *CommonCSVWriter) SetDestination(destination Destination) {
	writer.destination = destination
}

func (writer *CommonCSVWriter) GetHeaders() (headers Headers) {
	return writer.Headers
}

func (writer *CommonCSVWriter) GetTotalRows() (rows int) {
	if writer == nil {
		return
	}

	return writer.rows
}

func (writer *CommonCSVWriter) SetHeaders(headers Headers) {
	if writer == nil {
		return
	}

	writer.Headers = headers
}

func (writer *CommonCSVWriter) WriteHeaders() (err error) {
	if writer == nil {
		return
	}

	if writer.IsHeaderWritten() {
		return
	}

	if writer.shouldAddBOM {
		_, err = writer.file.Write([]byte{0xEF, 0xBB, 0xBF})
		if err != nil {
			log.Errorf("Error writing BOM to CSV file: %+v", err)
			return err
		}
	}

	// Write the header row.
	err = writer.Writer.Write(writer.GetHeaders())
	if err != nil {
		log.Errorf("Error writing header row to CSV file: %v", err)
		return err
	}

	writer.SetHeaderWritten()
	return

}

func (writer *CommonCSVWriter) WriteNextBatch(records Records) (err error) {
	if writer == nil {
		return fmt.Errorf("nil pointer dereference")
	}

	writer.rows += len(records)

	if len(records) == 0 {
		return
	}

	if !writer.IsHeaderWritten() {
		err = writer.WriteHeaders()
		if err != nil {
			return
		}
	}
	err = writer.Writer.WriteAll(records)
	if err != nil {
		log.Errorf("Error writing data rows to CSV file: %v", err)
		return err
	}

	// Check if any error occurred during the Flush operation.
	if err = writer.Writer.Error(); err != nil {
		log.Errorf("Error flushing CSV writer: %v", err)
		return
	}

	return
}

func (writer *CommonCSVWriter) Close(ctx context.Context) (err error) {
	err = writer.file.Close()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	switch writer.destination.GetType() {
	case DestinationLocal:
		return
	case DestinationS3:
		if writer.rows != 0 {
			err = writer.S3Upload(ctx)
			return
		}

		err = os.Remove(writer.destination.GetFileName())
		if err != nil {
			log.Errorf(err.Error())
			return
		}
	default:
		return
	}

	return

}

func (writer *CommonCSVWriter) GetUploadKey() string {
	return fmt.Sprintf(writer.destination.GetUploadDirectory() + writer.destination.GetFileName())
}

func (writer *CommonCSVWriter) S3Upload(ctx context.Context) (err error) {
	logTag := fmt.Sprintf(env.GetRequestID(ctx) + "function : Upload S3")

	file, err := os.Open(writer.destination.GetFileName())
	if err != nil {
		log.Errorf("Failed to open file, err:", err.Error())
	}
	defer file.Close()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.GetString(ctx, "s3.region")),
	})
	if err != nil {
		log.Errorf(logTag+"Could not create Session, err :", err.Error())
		return
	}

	// Create an S3 client
	uploader := s3manager.NewUploader(sess)

	// Upload the CSV fileName to S3

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(writer.destination.GetBucket()),
		Key:    aws.String(writer.GetUploadKey()),
		Body:   file,
	})

	if err != nil {
		log.Errorf(logTag+"Could not upload to S3, err :", err.Error())
		return
	}

	// Removing File from local After uploading
	err = os.Remove(writer.destination.GetFileName())
	if err != nil {
		log.Errorf("Failed to open file, err:", err.Error())
	}

	return nil
}

func (writer *CommonCSVWriter) GetPublicURL(ctx context.Context) (url string, err error) {
	logTag := fmt.Sprintf(env.GetRequestID(ctx) + "function : GetPresignedUrl ")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.GetString(ctx, "s3.region")),
	})

	if err != nil {
		log.Errorf(logTag+"Could not create Session, err :", err.Error())
		return
	}

	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(writer.destination.GetBucket()),
		Key:    aws.String(writer.GetUploadKey()),
	})

	url, err = req.Presign(2 * time.Hour)
	if err != nil {
		log.Errorf(logTag+"Failed to sign request, err :", err.Error())
		return
	}

	log.Infof(logTag+" The Presigned URL is %s", url)

	return url, nil
}

func (writer *CommonCSVWriter) Initialize() (err error) {
	csvFile, err := os.OpenFile(writer.destination.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Errorf("Unable to create or open CSV file: %v", err)
		return err
	}

	writer.file = csvFile
	writer.Writer = csv.NewWriter(csvFile)
	return
}
