package core

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	config "github.com/tommzn/go-config"
)

// newAWSConfig will create a new config for AWS.
// It try to find AWS region in passed config or in environment variable AWS_REGION.
func newAWSConfig(conf config.Config) *aws.Config {

	awsConfig := aws.NewConfig()
	if conf != nil {
		if awsRegion := conf.Get("aws.region", nil); awsRegion != nil {
			awsConfig.WithRegion(*awsRegion)
		}
	}

	if awsRegion, ok := os.LookupEnv("AWS_REGION"); ok {
		awsConfig.WithRegion(awsRegion)
	}
	return awsConfig
}
