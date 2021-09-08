package core

import (
	sqs "github.com/tommzn/aws-sqs"
	log "github.com/tommzn/go-log"
)

// A ScheduledCollector calls fetch method of a datasource one time and publishes returned event to a given AWS SQS queue.
// It contains a logger to provide insights to all processing steps and it requires a datasource and a puslisher for AWS SQS.
type ScheduledCollector struct {

	// Logger logs meesages and errors to a given output or log collector.
	logger log.Logger

	// Publisher sends events obtained from current datasource to defined AWS SQS queue.
	publisher sqs.Publisher

	// Queue defines the AWS SQS queue event from current datasource should be send to.
	queue string

	// ArchiveQueue is a queue all events are send additionally to.
	archiveQueue string

	// A datasource which fetches new data.
	datasource DataSource
}
