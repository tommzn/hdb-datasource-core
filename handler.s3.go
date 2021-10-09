package core

import (
	"context"

	awsevents "github.com/aws/aws-lambda-go/events"
	sqs "github.com/tommzn/aws-sqs"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

// NewS3EventHandler returns a new handler to process S3 events send from Cloud Watch.
func NewS3EventHandler(queue string, processor S3EventProcessor, conf config.Config, logger log.Logger) S3EventHandler {

	archiveQueue := conf.Get("hdb.archive", config.AsStringPtr("de.tsl.hdb.archive"))
	return &EventHandlerS3{
		logger:       logger,
		publisher:    sqs.NewPublisher(conf),
		queue:        queue,
		archiveQueue: *archiveQueue,
		processor:    processor,
	}
}

// Handle processes passed S3 event.
func (handler *EventHandlerS3) Handle(ctx context.Context, event awsevents.S3Event) error {
	return nil
}
