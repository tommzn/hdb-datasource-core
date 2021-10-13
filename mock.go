package core

import (
	"errors"

	awsevents "github.com/aws/aws-lambda-go/events"
	"github.com/golang/protobuf/proto"
	utils "github.com/tommzn/go-utils"
	events "github.com/tommzn/hdb-events-go"
)

// Mocked datasource which retruns a dummy event.
type dataSourceMock struct {
	shouldReturnError bool
}

// Fetch returns a dummy event for testing.
func (mock *dataSourceMock) Fetch() (proto.Message, error) {
	if mock.shouldReturnError {
		return nil, errors.New("Unable to fetch data.")
	} else {
		return &events.Dummy{Id: utils.NewId()}, nil
	}
}

// newDataSourceMock creates a new mocked datasource which returns always a new Dummy event.
func newDataSourceMock() DataSource {
	return &dataSourceMock{shouldReturnError: false}
}

// newDataSourceMockWithError creates a new mocked datasource which returns always an error.
func newDataSourceMockWithError() DataSource {
	return &dataSourceMock{shouldReturnError: true}
}

// publisherMock can be used to mock AWS SQS publisher for testing.
type publisherMock struct {
	shouldFail   bool
	messageCount int
}

// Counts calls to send message methods and returns a new message id. If you pass "error" as queue name it will returns with
// an error and doesn't count this call.
func (mock *publisherMock) send(message proto.Message) error {

	if mock.shouldFail {
		return errors.New("Unable to send message.")
	}

	mock.messageCount += 2
	return nil
}

// newPublisherMock returns a new mock for a AWS SQS publisher.
func newPublisherMock() publisher {
	return &publisherMock{shouldFail: false, messageCount: 0}
}

// s3EventProcessorMock can be used for testing, processing will aalways succeed.
type s3EventProcessorMock struct {
}

func (mock *s3EventProcessorMock) ProcessEvent(entity awsevents.S3Entity, content []byte) (proto.Message, error) {
	return nil, nil
}

func newS3EventProcessorMock() S3EventProcessor {
	return &s3EventProcessorMock{}
}
