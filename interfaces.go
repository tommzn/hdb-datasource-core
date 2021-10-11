package core

import (
	"context"

	awsevents "github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/protobuf/proto"
)

// DataSource retrieves data from a specific source.
type DataSource interface {

	// Fetch will retrieve new data and returns it as an event from central event lib.
	// For more details about events see https://github.com/tommzn/hdb-events-go
	Fetch() (proto.Message, error)
}

// A Collector calls fetch method of a datasource and process the returned event.
type Collector interface {

	// Run executes collectors data processing logic.
	Run() error
}

// SqsEventProcessor is used to handle event forwarded from AWS SQS to a lambda function.
type SqsEventProcessor interface {

	// Handle processes given SQS events.
	Handle(ctx context.Context, sqsEvent events.SQSEvent) error
}

// S3EventHandler is used to process an event published for S3 actions.
type S3EventHandler interface {

	// Handle processes passed S3 event.
	Handle(ctx context.Context, event awsevents.S3Event) error
}

// S3EventProcessor processes an event for a specific S3 object.
type S3EventProcessor interface {

	// Process is called to process given event for a S3 object.
	// It's always called, independent of return value from DownloadS3Object.
	ProcessEvent(event awsevents.S3Event) (proto.Message, error)

	// ProcessContent is called to process given content of a S3 object.
	// It's only called if DownloadS3Object returns true.
	ProcessContent(entity awsevents.S3Entity, content []byte) (proto.Message, error)

	// DownloadS3Object tells the S3 event handler if it should fownload content for a S3 object to process it.
	DownloadS3Object() bool
}

// s3Downloader returns content of given object in a S3 bucket.
type s3Downloader interface {

	// getObjectContent will download and return content for passed object in an S3 bucket.
	getObjectContent(bucket, key string) ([]byte, error)
}
