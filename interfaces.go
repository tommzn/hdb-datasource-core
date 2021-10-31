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
	Run(context.Context) error
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
	// If download option is enable via config it will pass S3 object content as well.
	ProcessEvent(entity awsevents.S3Entity, content []byte) (proto.Message, error)
}

// Publisher is used to send messages to one or multiple queues.
type Publisher interface {

	// Send will publish passed message to given queues.
	Send(message proto.Message) error
}
