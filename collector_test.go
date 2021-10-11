package core

import (
	"testing"

	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

type CollectorTestSuite struct {
	suite.Suite
	conf config.Config
}

func TestCollectorTestSuite(t *testing.T) {
	suite.Run(t, new(CollectorTestSuite))
}

func (suite *CollectorTestSuite) SetupTest() {
	suite.conf = loadConfigForTest(nil)
}

func (suite *CollectorTestSuite) TestCreateCollector() {

	collector := NewScheduledCollector("test-queue", newDataSourceMock(), suite.conf, log.NewLoggerFromConfig(suite.conf, nil))
	suite.NotNil(collector)
}

func (suite *CollectorTestSuite) TestFetchData() {

	collector := collectorForTest(suite.conf)

	err := collector.Run()
	suite.Nil(err)
	suite.Equal(2, collector.(*ScheduledCollector).publisher.(*publisherMock).messageCount)
}

func (suite *CollectorTestSuite) TestFetchDataWithError() {

	collector := collectorWithDataSourceErrorForTest(suite.conf)

	err := collector.Run()
	suite.NotNil(err)
	suite.Equal(0, collector.(*ScheduledCollector).publisher.(*publisherMock).messageCount)
}

func (suite *CollectorTestSuite) TestWithPublisherError() {

	collector := collectorForTest(suite.conf)
	collector.(*ScheduledCollector).queue = "error"

	err := collector.Run()
	suite.NotNil(err)
	suite.Equal(0, collector.(*ScheduledCollector).publisher.(*publisherMock).messageCount)
}

func (suite *CollectorTestSuite) TestWithArchivePublisherError() {

	collector := collectorForTest(suite.conf)
	collector.(*ScheduledCollector).archiveQueue = "error"

	err := collector.Run()
	suite.Nil(err)
	suite.Equal(1, collector.(*ScheduledCollector).publisher.(*publisherMock).messageCount)
}
