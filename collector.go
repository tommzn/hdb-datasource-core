package core

import (
	"encoding/base64"

	"github.com/golang/protobuf/proto"
	sqs "github.com/tommzn/aws-sqs"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

// NewScheduledCollector returns a new scheduled collector for given config.
func NewScheduledCollector(queue string, datasource DataSource, conf config.Config, logger log.Logger) Collector {

	archiveQueue := conf.Get("hdb.archive", config.AsStringPtr("de.tsl.hdb.archive"))
	return &ScheduledCollector{
		logger:       logger,
		publisher:    sqs.NewPublisher(conf),
		queue:        queue,
		archiveQueue: *archiveQueue,
		datasource:   datasource,
	}
}

// Run calls fetch of current datasource one time and published the returned event to a given AWS SQS queue.
// In can of any errors, they'll be logged and returned.
func (collector *ScheduledCollector) Run() error {

	defer collector.logger.Flush()

	event, err := collector.datasource.Fetch()
	if err != nil {
		collector.logger.Error("Unable to fetch new data, reason: ", err)
		return err
	}

	eventData, err := serializeEvent(event)
	if err != nil {
		collector.logger.Errorf("Failed to encode event, type: %T, reason: %s", event, err)
		return err
	}

	messageId, err := collector.publisher.Send(eventData, collector.queue)
	if err != nil {
		collector.logger.Error("Unable to semd event, reason: ", err)
		return err
	}
	collector.logger.Infof("Event send, type: %T, queue: %s, id: %s", event, collector.queue, *messageId)

	archiveMessageId, err := collector.publisher.SendAttributedMessage(eventData, collector.archiveQueue, map[string]string{ORIGIN_QUEUE: collector.queue})
	if err == nil {
		collector.logger.Info("Event send to archive queue, id: ", *archiveMessageId)
	} else {
		collector.logger.Error("Unable to semd event to archive queue, reason: ", err)
	}
	return nil
}

// serializeEvent uses protobuf to marshal given event
func serializeEvent(event proto.Message) (string, error) {

	protoContent, err := proto.Marshal(event)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(protoContent), nil
}
