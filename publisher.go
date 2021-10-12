package core

import (
	"github.com/golang/protobuf/proto"
	sqs "github.com/tommzn/aws-sqs"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

// newSqsPublisher creates a new SQS message publisher.
func newSqsPublisher(conf config.Config, logger log.Logger, queue, archiveQueue string) publisher {
	return &sqsPublisher{
		logger:       logger,
		sqsClient:    sqs.NewPublisher(conf),
		queue:        queue,
		archiveQueue: archiveQueue,
	}
}

// send will publish passed message to given queues.
func (publisher *sqsPublisher) send(message proto.Message) error {

	messageString, err := serializeEvent(message)
	if err != nil {
		publisher.logger.Errorf("Failed to encode event, type: %T, reason: %s", message, err)
		return err
	}

	messageId, err := publisher.sqsClient.Send(messageString, publisher.queue)
	if err != nil {
		publisher.logger.Error("Unable to semd event, reason: ", err)
		return err
	}
	publisher.logger.Infof("Event send, type: %T, queue: %s, id: %s", message, publisher.queue, *messageId)

	archiveMessageId, err := publisher.sqsClient.SendAttributedMessage(messageString, publisher.archiveQueue, map[string]string{ORIGIN_QUEUE: publisher.queue})
	if err != nil {
		publisher.logger.Error("Unable to semd event to archive queue, reason: ", err)
		return err
	}
	publisher.logger.Info("Event send to archive queue, id: ", *archiveMessageId)

	return nil
}
