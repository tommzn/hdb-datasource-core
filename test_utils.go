package core

import (
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

// loadConfigForTest loads test config.
func loadConfigForTest(fileName *string) config.Config {

	configFile := "testconfig.yml"
	if fileName != nil {
		configFile = *fileName
	}
	configLoader := config.NewFileConfigSource(&configFile)
	config, _ := configLoader.Load()
	return config
}

// collectorForTest returns a collector for testing which uses a dummy datasource and a mocked AWS SQS publisher.
func collectorForTest(conf config.Config) Collector {
	return &ScheduledCollector{
		logger:       log.NewLoggerFromConfig(conf, nil),
		publisher:    newPublisherMock(),
		queue:        "test-queue",
		archiveQueue: "test-archive",
		datasource:   newDataSourceMock(),
	}
}

// collectorWithDataSourceErrorForTest returns a collector for testing which uses a always failing datasource
// and a mocked AWS SQS publisher.
func collectorWithDataSourceErrorForTest(conf config.Config) Collector {
	return &ScheduledCollector{
		logger:       log.NewLoggerFromConfig(conf, nil),
		publisher:    newPublisherMock(),
		queue:        "test-queue",
		archiveQueue: "test-archive",
		datasource:   newDataSourceMockWithError(),
	}
}
