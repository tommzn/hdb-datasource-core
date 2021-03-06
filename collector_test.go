package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

type CollectorTestSuite struct {
	suite.Suite
	conf config.Config
	ctx  context.Context
}

func TestCollectorTestSuite(t *testing.T) {
	suite.Run(t, new(CollectorTestSuite))
}

func (suite *CollectorTestSuite) SetupTest() {
	suite.conf = loadConfigForTest(nil)
	suite.ctx = context.Background()
}

func (suite *CollectorTestSuite) TestCreateCollector() {

	collector := NewScheduledCollector("test-queue", newDataSourceMock(), suite.conf, log.NewLoggerFromConfig(suite.conf, nil))
	suite.NotNil(collector)
}

func (suite *CollectorTestSuite) TestFetchData() {

	collector := collectorForTest(suite.conf)

	err := collector.Run(suite.ctx)
	suite.Nil(err)
	suite.Equal(2, collector.(*ScheduledCollector).messagePublisher.(*publisherMock).messageCount)
}

func (suite *CollectorTestSuite) TestFetchDataWithError() {

	collector := collectorWithDataSourceErrorForTest(suite.conf)

	err := collector.Run(suite.ctx)
	suite.NotNil(err)
	suite.Equal(0, collector.(*ScheduledCollector).messagePublisher.(*publisherMock).messageCount)
}

func (suite *CollectorTestSuite) TestWithPublisherError() {

	collector := collectorForTest(suite.conf)
	collector.(*ScheduledCollector).messagePublisher.(*publisherMock).shouldFail = true

	err := collector.Run(suite.ctx)
	suite.NotNil(err)
	suite.Equal(0, collector.(*ScheduledCollector).messagePublisher.(*publisherMock).messageCount)
}

func (suite *CollectorTestSuite) TestWithArchivePublisherError() {

	collector := collectorForTest(suite.conf)
	collector.(*ScheduledCollector).messagePublisher.(*publisherMock).shouldFail = true

	err := collector.Run(suite.ctx)
	suite.NotNil(err)
	suite.Equal(0, collector.(*ScheduledCollector).messagePublisher.(*publisherMock).messageCount)
}

func (suite *CollectorTestSuite) TestContinuousCollector() {

	mock := newCollectorMock()
	collector := NewContinuousCollector(mock, loggerForTest())
	collector.(*ContinuousCollector).signalObserver = osSignalObserverMock

	ctx := context.Background()
	collector.Run(ctx)

	suite.True(mock.(*collectorMock).hasFinished)
	suite.True(mock.(*collectorMock).counter > 0)
}
