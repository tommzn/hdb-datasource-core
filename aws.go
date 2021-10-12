package core

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	config "github.com/tommzn/go-config"
)

// newAWSConfig try to find AWS region in passed config or in environment variable AWS_REGION
// and returns a new AWS config.
func newAWSConfig(conf config.Config) *aws.Config {

	awsConfig := aws.NewConfig()
	if awsRegion, ok := os.LookupEnv("AWS_REGION"); ok {
		return awsConfig.WithRegion(awsRegion)
	}

	if conf != nil {
		configKeys := []string{"aws.region", "aws.s3.region"}
		for _, configKey := range configKeys {
			if awsRegion := conf.Get(configKey, nil); awsRegion != nil {
				return awsConfig.WithRegion(*awsRegion)
			}
		}

	}
	return awsConfig
}

func newS3Downloader(conf config.Config) *s3manager.Downloader {
	return s3manager.NewDownloader(newAwsSession(conf))
}

func newAwsSession(conf config.Config) *session.Session {
	return session.Must(session.NewSession(newAWSConfig(conf)))
}
