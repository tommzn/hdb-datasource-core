package core

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"google.golang.org/protobuf/proto"

	awsevents "github.com/aws/aws-lambda-go/events"
	goconfig "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

// NewS3EventHandler returns a new handler to process S3 events sent from CloudWatch.
func NewS3EventHandler(queue string, processor S3EventProcessor, conf goconfig.Config, logger log.Logger) S3EventHandler {

	downloadContent := conf.GetAsBool("aws.s3.download", goconfig.AsBoolPtr(false))
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

	defer handler.logger.Flush()

	var errorList []error
	for _, record := range event.Records {

		message, err := handler.processS3Entity(ctx, record.S3)
		if err != nil {
			handler.logger.Errorf("S3 event processing failed %s/%s, reason: %s", record.S3.Bucket.Name, record.S3.Object.Key, err)
			errorList = append(errorList, err)
			continue
		}

		if err := handler.messagePublisher.Send(message); err != nil {
			handler.logger.Errorf("Unable to send event for S3 entity %s/%s, reason: %s", record.S3.Bucket.Name, record.S3.Object.Key, err)
			errorList = append(errorList, err)
		}
	}
	return asError(errorList)
}

// processS3Entity processes a single S3 entity.
func (handler *EventHandlerS3) processS3Entity(ctx context.Context, entity awsevents.S3Entity) (proto.Message, error) {

	content, err := handler.downloadS3ObjectIfNecessary(ctx, entity)
	if err != nil {
		return nil, err
	}
	return handler.processor.ProcessEvent(entity, content)
}

// downloadS3ObjectIfNecessary will return S3 object content if processor requires it.
func (handler *EventHandlerS3) downloadS3ObjectIfNecessary(ctx context.Context, entity awsevents.S3Entity) ([]byte, error) {

	if handler.downloader != nil {
		content, err := handler.getObjectContent(ctx, entity)
		if err != nil {
			handler.logger.Errorf("Unable to download %s/%s, reason: %s", entity.Bucket.Name, entity.Object.Key, err)
		}
		return content, err
	}
	return []byte{}, nil
}

// getObjectContent tries to download content for given S3 object.
func (handler *EventHandlerS3) getObjectContent(ctx context.Context, entity awsevents.S3Entity) ([]byte, error) {

	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(entity.Bucket.Name),
		Key:    aws.String(entity.Object.Key),
	}

	buf := manager.NewWriteAtBuffer([]byte{})
	_, err := handler.downloader.Download(ctx, buf, getObjectInput)
	if err != nil {
		return nil, err
	}
	// Convert to []byte
	return io.ReadAll(bytes.NewReader(buf.Bytes()))
}
