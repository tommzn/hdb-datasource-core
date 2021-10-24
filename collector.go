package core

import (
	"context"
	"encoding/base64"

	"github.com/golang/protobuf/proto"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

// NewScheduledCollector returns a new scheduled collector for given config.
func NewScheduledCollector(queue string, datasource DataSource, conf config.Config, logger log.Logger) Collector {

	return &ScheduledCollector{
		logger:           logger,
		messagePublisher: newSqsPublisher(conf, logger, queue, archiveQueueFromConfig(conf)),
		datasource:       datasource,
	}
}

// Run calls fetch of current datasource one time and published the returned event to a given AWS SQS queue.
// In can of any errors, they'll be logged and returned.
func (collector *ScheduledCollector) Run(ctx context.Context) error {

	defer collector.logger.Flush()

	event, err := collector.datasource.Fetch()
	if err != nil {
		collector.logger.Error("Unable to fetch new data, reason: ", err)
		return err
	}
	return collector.messagePublisher.send(event)
}

// serializeEvent uses protobuf to marshal given event
func serializeEvent(event proto.Message) (string, error) {

	protoContent, err := proto.Marshal(event)
	return base64.StdEncoding.EncodeToString(protoContent), err
}
