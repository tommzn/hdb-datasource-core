package core

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	goconfig "github.com/tommzn/go-config"
)

// newAWSConfig tries to find AWS region in passed config or in environment variable AWS_REGION
// and returns a new AWS v2 config.
func newAWSConfig(conf goconfig.Config) aws.Config {

	// Try config keys first
	if conf != nil {
		configKeys := []string{"aws.region", "aws.s3.region"}
		for _, configKey := range configKeys {
			if awsRegion := conf.Get(configKey, nil); awsRegion != nil {
				cfg, _ := config.LoadDefaultConfig(context.Background(), config.WithRegion(*awsRegion))
				return cfg
			}
		}
	}

	// Fallback to environment
	if awsRegion, ok := os.LookupEnv("AWS_REGION"); ok {
		cfg, _ := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
		return cfg
	}

	cfg, _ := config.LoadDefaultConfig(context.Background())
	return cfg
}

// newS3Downloader creates an AWS SDK v2 S3 downloader.
func newS3Downloader(conf goconfig.Config) *manager.Downloader {
	awsCfg := newAWSConfig(conf)
	client := s3.NewFromConfig(awsCfg)
	return manager.NewDownloader(client)
}
