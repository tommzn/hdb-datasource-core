package core

import (
	"github.com/golang/protobuf/proto"
	sqs "github.com/tommzn/aws-sqs"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

// newSqsPublisher creates a new SQS message publisher.
func NewPublisher(conf config.Config, logger log.Logger) Publisher {
	queue := conf.Get("hdb.queue", config.AsStringPtr("de.tsl.hdb.unknown"))
	archiveQueue := archiveQueueFromConfig(conf)
	return newSqsPublisher(conf, logger, *queue, archiveQueue)
}

// newSqsPublisher creates a new SQS message publisher with given queue and archive queue.
func newSqsPublisher(conf config.Config, logger log.Logger, queue, archiveQueue string) Publisher {
	return &SqsPublisher{
		logger:       logger,
		sqsClient:    sqs.NewPublisher(conf),
		queue:        queue,
		archiveQueue: archiveQueue,
	}
}

// send will publish passed message to given queues.
func (publisher *SqsPublisher) Send(message proto.Message) error {

	defer publisher.logger.Flush()
	logEvent(message, publisher.logger)

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
		publisher.logger.Errorf("Unable to semd event to archive queue %s, reason: %s", publisher.archiveQueue, err)
		return err
	}
	publisher.logger.Info("Event send to archive queue, id: ", *archiveMessageId)

	return nil
}
