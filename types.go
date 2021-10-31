package core

import (
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	sqs "github.com/tommzn/aws-sqs"
	log "github.com/tommzn/go-log"
)

const (
	// ORIGIN_QUEUE is used to add name of a source queue to message attributes of archive events.
	ORIGIN_QUEUE string = "origin_queue"
)

// A ScheduledCollector calls fetch method of a datasource one time and publishes returned event to a given AWS SQS queue.
// It contains a logger to provide insights to all processing steps and it requires a datasource and a puslisher for AWS SQS.
type ScheduledCollector struct {

	// Logger logs meesages and errors to a given output or log collector.
	logger log.Logger

	// Publisher sends events obtained from current datasource to defined AWS SQS queue.
	messagePublisher Publisher

	// A datasource which fetches new data.
	datasource DataSource
}

// ContinuousCollector is used as a deamon to permanently collect data from a source.
// It mainly cares about observing os singles to handle graceful shutdowns. The actual
// logic to process data in encapsulated in datasource member.
type ContinuousCollector struct {

	// Logger logs meesages and errors to a given output or log collector.
	logger log.Logger

	// Datasource is a collector as well which contains the actual logic to process data from
	// different sources.
	datasource Collector

	// signalObserver waits for OS signals like SIGINT and SIGTERM.
	signalObserver osSignalObserver
}

// EventHandlerS3 is used to process an S3 event send from Cloud Watch to a Lambda function on AWS.
type EventHandlerS3 struct {

	// Logger logs meesages and errors to a given output or log collector.
	logger log.Logger

	// Publisher sends events obtained from current datasource to defined AWS SQS queue.
	messagePublisher Publisher

	// processor will be called to process a received event.
	processor S3EventProcessor

	// downloader is used to get objet content for an object in a S3 bucket.
	downloader *s3manager.Downloader
}

// SqsPublisher is used to publish messages on AWS SQS.
type SqsPublisher struct {

	// sqsClient sends events obtained from current datasource to defined AWS SQS queue.
	sqsClient sqs.Publisher

	// Logger logs meesages and errors to a given output or log collector.
	logger log.Logger

	// Queue defines the AWS SQS queue event from current datasource should be send to.
	queue string

	// ArchiveQueue is a queue all events are send additionally to.
	archiveQueue string
}

// osSignalObserver will observer OS signals. Execution is blocked until a signal has been received.
type osSignalObserver = func()
