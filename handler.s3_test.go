package core

import (
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
