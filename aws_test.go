package core

import (
	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
	"os"
	"testing"
)

type AwsTestSuite struct {
	suite.Suite
}

func TestAwsTestSuite(t *testing.T) {
	suite.Run(t, new(AwsTestSuite))
}

func (suite *AwsTestSuite) TestAwsConfig() {

	awsConfig1 := newAWSConfig(nil)
	suite.NotNil(awsConfig1)
	suite.Nil(awsConfig1.Region)

	conf := loadConfigForTest(config.AsStringPtr("fixtures/aws.testconfig.yml"))
	awsConfig2 := newAWSConfig(conf)
	suite.NotNil(awsConfig2)
	suite.NotNil(awsConfig2.Region)
	suite.Equal("eu-south-4", *awsConfig2.Region)

	expectedRegion := "eu-east-17"
	os.Setenv("AWS_REGION", expectedRegion)
	awsConfig3 := newAWSConfig(conf)
	suite.NotNil(awsConfig3)
	suite.NotNil(awsConfig3.Region)
	suite.Equal(expectedRegion, *awsConfig3.Region)
	os.Unsetenv("AWS_REGION")
}
