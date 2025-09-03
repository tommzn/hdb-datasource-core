package core

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/tommzn/go-log"
)

const (
	// ORIGIN_QUEUE is used to add name of a source queue to message attributes of archive events.
	ORIGIN_QUEUE string = "origin_queue"
)

// A ScheduledCollector calls fetch method of a datasource one time and publishes returned event to a given AWS SQS queue.
// It contains a logger to provide insights to all processing steps and it requires a datasource and a publisher for AWS SQS.
type ScheduledCollector struct {

	// Logger logs messages and errors to a given output or log collector.
	logger log.Logger

	// Publisher sends events obtained from current datasource to defined AWS SQS queue.
	messagePublisher Publisher

	// A datasource which fetches new data.
	datasource DataSource
}

// ContinuousCollector is used as a daemon to permanently collect data from a source.
// It mainly cares about observing OS signals to handle graceful shutdowns. The actual
// logic to process data is encapsulated in datasource member.
type ContinuousCollector struct {

	// Logger logs messages and errors to a given output or log collector.
	logger log.Logger

	// Datasource is a collector as well which contains the actual logic to process data from
	// different sources.
	datasource Collector

	// signalObserver waits for OS signals like SIGINT and SIGTERM.
	signalObserver osSignalObserver
}

// EventHandlerS3 is used to process an S3 event sent from CloudWatch to a Lambda function on AWS.
type EventHandlerS3 struct {

	// Logger logs messages and errors to a given output or log collector.
	logger log.Logger

	// Publisher sends events obtained from current datasource to defined AWS SQS queue.
	messagePublisher Publisher

	// Processor will be called to process a received event.
	processor S3EventProcessor

	// Downloader is used to get object content for an object in an S3 bucket.
	downloader *s3.Downloader
}

// SqsPublisher is used to publish messages on AWS SQS.
type SqsPublisher struct {

	// sqsClient sends events obtained from current datasource to defined AWS SQS queue.
	sqsClient *sqs.Client

	// Logger logs messages and errors to a given output or log collector.
	logger log.Logger

	// Queue defines the AWS SQS queue event from current datasource should be sent to.
	queueURL string

	// ArchiveQueue is a queue all events are sent additionally to.
	archiveQueue string
}

// osSignalObserver will observe OS signals. Execution is blocked until a signal has been received.
type osSignalObserver = func()
