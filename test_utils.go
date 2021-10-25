package core

import (
	"os"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

// loadConfigForTest loads test config.
func loadConfigForTest(fileName *string) config.Config {

	configFile := "fixtures/testconfig.yml"
	if fileName != nil {
		configFile = *fileName
	}
	configLoader := config.NewFileConfigSource(&configFile)
	config, _ := configLoader.Load()
	return config
}

// loggerForTest creates a new stdout logger for testing.
func loggerForTest() log.Logger {
	return log.NewLogger(log.Debug, nil, nil)
}

// collectorForTest returns a collector for testing which uses a dummy datasource and a mocked AWS SQS publisher.
func collectorForTest(conf config.Config) Collector {
	return &ScheduledCollector{
		logger:           loggerForTest(),
		messagePublisher: newPublisherMock(),
		datasource:       newDataSourceMock(),
	}
}

// collectorWithDataSourceErrorForTest returns a collector for testing which uses a always failing datasource
// and a mocked AWS SQS publisher.
func collectorWithDataSourceErrorForTest(conf config.Config) Collector {
	return &ScheduledCollector{
		logger:           loggerForTest(),
		messagePublisher: newPublisherMock(),
		datasource:       newDataSourceMockWithError(),
	}
}

func newS3EventHandlerForTest(conf config.Config) S3EventHandler {
	return &EventHandlerS3{
		logger:           loggerForTest(),
		messagePublisher: newPublisherMock(),
		processor:        newS3EventProcessorMock(),
		downloader:       newS3Downloader(conf),
	}
}

func s3EntityForTest() events.S3Entity {
	return events.S3Entity{
		Bucket: events.S3Bucket{
			Name: os.Getenv("AWS_S3_TEST_BUCKET"),
		},
		Object: events.S3Object{
			Key: "s3.download.test",
		},
	}
}

func notExistingS3EntityForTest() events.S3Entity {
	return events.S3Entity{
		Bucket: events.S3Bucket{
			Name: os.Getenv("AWS_S3_TEST_BUCKET"),
		},
		Object: events.S3Object{
			Key: "s3.failue.test",
		},
	}
}

func s3EventForTest(entity events.S3Entity) events.S3Event {
	return events.S3Event{
		Records: []events.S3EventRecord{
			events.S3EventRecord{
				S3: entity,
			},
		},
	}
}

// osSignalObserverMock can be used to simulate observing for OS signals.
// It will return after 1 second.
func osSignalObserverMock() {
	time.Sleep(1 * time.Second)
}

// skipCI returns true if env variable CI is set
func skipCI(t *testing.T) {
	if _, isSet := os.LookupEnv("CI"); isSet {
		t.Skip("Skipping testing in CI environment")
	}
}
