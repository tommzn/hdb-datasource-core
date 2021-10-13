package core

import (
	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
	events "github.com/tommzn/hdb-events-go"
	"testing"
)

type PublisherTestSuite struct {
	suite.Suite
	conf config.Config
}

func TestPublisherTestSuite(t *testing.T) {
	suite.Run(t, new(PublisherTestSuite))
}

func (suite *PublisherTestSuite) SetupTest() {
	suite.conf = loadConfigForTest(nil)
}

func (suite *PublisherTestSuite) TestPublishMessage() {

	skipCI(suite.T())

	publisher := newSqsPublisher(suite.conf, loggerForTest(), "tzn-unittest", "tzn-unittest")
	suite.NotNil(publisher)

	event := &events.Dummy{Id: "1"}
	suite.Nil(publisher.send(event))
}

func (suite *PublisherTestSuite) TestPublishWithError() {

	skipCI(suite.T())

	event := &events.Dummy{Id: "1"}

	publisher1 := newSqsPublisher(suite.conf, loggerForTest(), "xxx", "tzn-unittest")
	suite.NotNil(publisher1)
	suite.NotNil(publisher1.send(event))

	publisher2 := newSqsPublisher(suite.conf, loggerForTest(), "tzn-unittest", "xxx")
	suite.NotNil(publisher2)
	suite.NotNil(publisher2.send(event))
}
