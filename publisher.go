package core

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/tommzn/go-config"
	"github.com/tommzn/go-log"
	"google.golang.org/protobuf/proto"
)

// NewPublisher creates a new SQS message publisher.
func NewPublisher(conf config.Config, logger log.Logger) Publisher {
	queue := conf.Get("hdb.queue", config.AsStringPtr("de.tsl.hdb.unknown"))
	archiveQueue := archiveQueueFromConfig(conf)
	return newSqsPublisher(conf, logger, *queue, archiveQueue)
}

// newSqsPublisher creates a new SQS message publisher with given queue and archive queue.
func newSqsPublisher(conf config.Config, logger log.Logger, queue, archiveQueue string) Publisher {

	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}

	return &SqsPublisher{
		logger:       logger,
		sqsClient:    sqs.NewFromConfig(awsCfg),
		queueURL:     queue,
		archiveQueue: archiveQueue,
	}
}

// Send will publish passed message to given queues.
func (publisher *SqsPublisher) Send(message proto.Message) error {

	defer publisher.logger.Flush()
	logEvent(message, publisher.logger)

	messageString, err := serializeEvent(message)
	if err != nil {
		publisher.logger.Errorf("Failed to encode event, type: %T, reason: %s", message, err)
		return err
	}

	// Send to primary queue
	sendOut, err := publisher.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(publisher.queueURL),
		MessageBody: aws.String(messageString),
	})
	if err != nil {
		publisher.logger.Error("Unable to send event, reason: ", err)
		return err
	}
	publisher.logger.Infof("Event sent, type: %T, queue: %s, id: %s", message, publisher.queueURL, *sendOut.MessageId)

	// Send to archive queue with attributes
	archiveOut, err := publisher.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(publisher.archiveQueue),
		MessageBody: aws.String(messageString),
		MessageAttributes: map[string]types.MessageAttributeValue{
			ORIGIN_QUEUE: {
				DataType:    aws.String("String"),
				StringValue: aws.String(publisher.queueURL),
			},
		},
	})
	if err != nil {
		publisher.logger.Errorf("Unable to send event to archive queue %s, reason: %s", publisher.archiveQueue, err)
		return err
	}
	publisher.logger.Info("Event sent to archive queue, id: ", *archiveOut.MessageId)

	return nil
}
