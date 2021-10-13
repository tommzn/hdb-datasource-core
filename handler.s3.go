package core

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/protobuf/proto"

	awsevents "github.com/aws/aws-lambda-go/events"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

// NewS3EventHandler returns a new handler to process S3 events send from Cloud Watch.
func NewS3EventHandler(queue string, processor S3EventProcessor, conf config.Config, logger log.Logger) S3EventHandler {

	downloadContent := conf.GetAsBool("aws.s3.download", config.AsBoolPtr(false))
	handler := &EventHandlerS3{
		logger:           logger,
		messagePublisher: newSqsPublisher(conf, logger, queue, archiveQueueFromConfig(conf)),
		processor:        processor,
	}
	if *downloadContent {
		handler.downloader = newS3Downloader(conf)
	}
	return handler
}

// Handle processes passed S3 event.
func (handler *EventHandlerS3) Handle(ctx context.Context, event awsevents.S3Event) error {

	var errorList []error
	for _, record := range event.Records {

		message, err := handler.processS3Entity(record.S3)
		if err != nil {
			handler.logger.Errorf("S3 event processing failed %s/%s, reason: ", record.S3.Bucket.Name, record.S3.Object.Key, err)
			errorList = append(errorList, err)
		}

		if err := handler.messagePublisher.send(message); err != nil {
			errorList = append(errorList, err)
		}
	}
	return asError(errorList)
}

// processS3Entity process a single S3 entity.
func (handler *EventHandlerS3) processS3Entity(entity awsevents.S3Entity) (proto.Message, error) {

	content, err := handler.downloadS3ObjectIfNecessary(entity)
	if err != nil {
		return nil, err
	}
	return handler.processor.ProcessEvent(entity, content)
}

// downloadS3ObjectIfNecessary will return S3 object content if processor returns
// true at DownloadS3Object.
func (handler *EventHandlerS3) downloadS3ObjectIfNecessary(entity awsevents.S3Entity) ([]byte, error) {

	if handler.downloader != nil {
		content, err := handler.getObjectContent(entity)
		if err != nil {
			handler.logger.Errorf("Unable to download %s/%s, reason: %s", entity.Bucket.Name, entity.Object.Key, err)
		}
		return content, err
	}
	return []byte{}, nil
}

// getObjectContent tries to download content for given S3 object.
func (handler *EventHandlerS3) getObjectContent(entity awsevents.S3Entity) ([]byte, error) {

	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(entity.Bucket.Name),
		Key:    aws.String(entity.Object.Key),
	}

	buf := &aws.WriteAtBuffer{}
	_, err := handler.downloader.Download(buf, getObjectInput)
	return buf.Bytes(), err
}
