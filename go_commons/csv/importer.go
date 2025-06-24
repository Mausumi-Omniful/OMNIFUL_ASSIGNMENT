package csv

import (
	"context"
	"errors"

	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/s3"
)

const nilPointerRef = "nil Pointer Reference"
const notificationDisabled = "notification are disabled for this import in client. Turn on while creating client"

type Importer struct {
	oReader               *CommonCSV
	oWriter               *CommonCSVWriter
	notifyFailedEntries   bool
	isWriterHeaderWritten bool
}

func (importer *Importer) GetReader() (csvReader *CommonCSV) {
	if importer == nil {
		return
	}

	return importer.oReader
}

func (importer *Importer) SetReader(reader *CommonCSV) {
	if importer == nil {
		return
	}

	importer.oReader = reader
}

func (importer *Importer) GetWriter() (csvWriter *CommonCSVWriter) {
	if importer == nil {
		return
	}

	return importer.oWriter
}

func (importer *Importer) SetWriter(csvWriter *CommonCSVWriter) {
	if importer == nil {
		return
	}

	importer.oWriter = csvWriter
}

func (importer *Importer) ShouldNotifyFailedEntries() (notify bool) {
	if importer == nil {
		return
	}

	return importer.notifyFailedEntries
}

func (importer *Importer) SetShouldNotifyFailedEntries(notify bool) {
	if importer == nil {
		return
	}

	importer.notifyFailedEntries = notify
}

func (importer *Importer) Initialize(ctx context.Context) (err error) {
	if importer == nil {
		return errors.New(nilPointerRef)
	}

	err = importer.oReader.InitializeReader(ctx)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	if importer.ShouldNotifyFailedEntries() {
		err = importer.oWriter.Initialize()
		if err != nil {
			log.Errorf(err.Error())
			return
		}
	}

	return
}

func NewCSVImporter(options ...ImporterOption) (*Importer, error) {
	opts := &ImporterOptions{}

	for _, option := range options {
		option(opts)
	}

	return &Importer{
		oReader:             opts.oReader,
		oWriter:             opts.oWriter,
		notifyFailedEntries: opts.notifyFailedEntries,
	}, nil
}

func (importer *Importer) IsEOF() (eof bool) {
	if importer == nil {
		return true
	}

	return importer.oReader.IsEOF()
}

func (importer *Importer) ParseNextBatch(data interface{}) (err error) {
	if importer == nil {
		return errors.New(nilPointerRef)
	}

	return importer.oReader.ParseNextBatch(data)
}

func (importer *Importer) setWriterHeaderWritten() {
	importer.isWriterHeaderWritten = true
}

func (importer *Importer) WriteRawHeaders(headers Headers) (err error) {
	if importer == nil {
		log.Errorf(nilPointerRef)
		return errors.New(nilPointerRef)
	}

	if !importer.ShouldNotifyFailedEntries() {
		log.Errorf(notificationDisabled)
		return errors.New(notificationDisabled)
	}

	if importer.GetWriter() == nil {
		log.Errorf(nilPointerRef)
		return errors.New(nilPointerRef)
	}

	importer.oWriter.SetHeaders(headers)

	return importer.oWriter.WriteHeaders()
}

func (importer *Importer) WriteNextRawBatch(records Records) (err error) {
	if importer == nil {
		log.Errorf(nilPointerRef)
		return errors.New(nilPointerRef)
	}

	if !importer.ShouldNotifyFailedEntries() {
		log.Errorf(notificationDisabled)
		return errors.New(notificationDisabled)
	}

	if importer.GetWriter() == nil {
		log.Errorf(nilPointerRef)
		return errors.New(nilPointerRef)
	}

	if !importer.isWriterHeaderWritten {
	}

	return importer.oWriter.WriteNextBatch(records)
}

func (importer *Importer) WriteNextBatch(dataProvider DataProvider) (err error) {
	if importer == nil {
		log.Errorf(nilPointerRef)
		return errors.New(nilPointerRef)
	}

	if !importer.ShouldNotifyFailedEntries() {
		log.Errorf(notificationDisabled)
		return errors.New(notificationDisabled)
	}

	if importer.GetWriter() == nil {
		log.Errorf(nilPointerRef)
		return errors.New(nilPointerRef)
	}

	if !importer.isWriterHeaderWritten {
		err = importer.WriteHeaders(dataProvider)
		if err != nil {
			return
		}
	}

	if dataProvider == nil {
		return
	}

	return importer.oWriter.WriteNextBatch(dataProvider.GetDataRows())
}

func (importer *Importer) GetPreSignedURL(ctx context.Context) (url string, err error) {
	if importer == nil {
		log.Errorf(nilPointerRef)
		err = errors.New(nilPointerRef)
		return
	}

	if !importer.ShouldNotifyFailedEntries() {
		log.Errorf(notificationDisabled)
		err = errors.New(notificationDisabled)
		return
	}

	if importer.GetWriter() == nil {
		log.Errorf(nilPointerRef)
		err = errors.New(nilPointerRef)
		return
	}

	return importer.oWriter.GetPublicURL(ctx)
}

func (importer *Importer) GetPublicBucketPermanentURL(ctx context.Context) (url string, err error) {
	if importer == nil {
		log.Errorf(nilPointerRef)
		err = errors.New(nilPointerRef)
		return
	}

	if !importer.ShouldNotifyFailedEntries() {
		log.Errorf(notificationDisabled)
		err = errors.New(notificationDisabled)
		return
	}

	if importer.GetWriter() == nil {
		log.Errorf(nilPointerRef)
		err = errors.New(nilPointerRef)
		return
	}

	return s3.GetURL(ctx, s3.FileConfig{
		Bucket: importer.oWriter.destination.GetBucket(),
		Path:   importer.oWriter.GetUploadKey(),
	}), nil
}

func (importer *Importer) WriteHeaders(dataProvider DataProvider) (err error) {
	if importer == nil {
		log.Errorf(nilPointerRef)
		return errors.New(nilPointerRef)
	}

	if !importer.ShouldNotifyFailedEntries() {
		log.Errorf(notificationDisabled)
		return errors.New(notificationDisabled)
	}

	if importer.GetWriter() == nil {
		log.Errorf(nilPointerRef)
		return errors.New(nilPointerRef)
	}

	if dataProvider == nil {
		return
	}

	importer.oWriter.SetHeaders(dataProvider.GetHeaderRow())

	return importer.oWriter.WriteHeaders()
}

func (importer *Importer) SetWriterHeaders(headers Headers) {
	if importer == nil {
		log.Errorf(nilPointerRef)
		return
	}

	if !importer.ShouldNotifyFailedEntries() {
		log.Errorf(notificationDisabled)
		return
	}

	if importer.GetWriter() == nil {
		log.Errorf(nilPointerRef)
		return
	}

	importer.oWriter.SetHeaders(headers)
}

func (importer *Importer) GetWriterTotalRows() (rows int) {
	if importer == nil {
		log.Errorf(nilPointerRef)
		return
	}

	if !importer.ShouldNotifyFailedEntries() {
		log.Errorf(notificationDisabled)
		return
	}

	if importer.GetWriter() == nil {
		log.Errorf(nilPointerRef)
		return
	}

	return importer.oWriter.GetTotalRows()
}

func (importer *Importer) Close(ctx context.Context) (err error) {
	if importer == nil {
		return errors.New(nilPointerRef)
	}

	if importer.ShouldNotifyFailedEntries() {
		err = importer.oWriter.Close(ctx)
		if err != nil {
			log.Errorf(err.Error())
			return
		}
	}

	return
}
