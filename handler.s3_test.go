package core

import (
	"context"
	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
	"testing"
)

type S3HandlerTestSuite struct {
	suite.Suite
	conf config.Config
}

func TestS3HandlerTestSuite(t *testing.T) {
	suite.Run(t, new(S3HandlerTestSuite))
}

func (suite *S3HandlerTestSuite) SetupTest() {
	suite.conf = loadConfigForTest(nil)
}

func (suite *S3HandlerTestSuite) TestCreateHandler() {

	suite.NotNil(NewS3EventHandler("test-queue", newS3EventProcessorMock(), suite.conf, loggerForTest()))
}

func (suite *S3HandlerTestSuite) TestDownloadObjectContent() {

	skipCI(suite.T())

	handler := newS3EventHandlerForTest(suite.conf)
	entity := s3EntityForTest()
	content1, err1 := handler.(*EventHandlerS3).getObjectContent(entity)
	suite.Nil(err1)
	suite.True(len(content1) > 0)

	entity.Object.Key = "xxx"
	content2, err2 := handler.(*EventHandlerS3).getObjectContent(entity)
	suite.NotNil(err2)
	suite.Len(content2, 0)

}

func (suite *S3HandlerTestSuite) TestProcessEvent() {

	skipCI(suite.T())

	handler := newS3EventHandlerForTest(suite.conf)

	event1 := s3EventForTest(s3EntityForTest())
	suite.Nil(handler.Handle(context.Background(), event1))

	event2 := s3EventForTest(notExistingS3EntityForTest())
	suite.NotNil(handler.Handle(context.Background(), event2))

	// Remove downloader to skip content download from S3
	handler.(*EventHandlerS3).downloader = nil
	suite.Nil(handler.Handle(context.Background(), event2))
}

func (suite *S3HandlerTestSuite) TestMessageSendFailure() {

	skipCI(suite.T())

	handler := newS3EventHandlerForTest(suite.conf)
	handler.(*EventHandlerS3).messagePublisher.(*publisherMock).shouldFail = true

	event1 := s3EventForTest(s3EntityForTest())
	suite.NotNil(handler.Handle(context.Background(), event1))
}
