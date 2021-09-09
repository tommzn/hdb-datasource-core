package core

import (
	"errors"

	"github.com/golang/protobuf/proto"
	sqs "github.com/tommzn/aws-sqs"
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
	messageCount int
}

// Mocked method to send messages to AWS SQS.
func (mock *publisherMock) Send(message interface{}, queue string) (*string, error) {
	return mock.send(queue)
}

// Mocked method to send attributed messages to AWS SQS.
func (mock *publisherMock) SendAttributedMessage(message interface{}, queue string, attributes map[string]string) (*string, error) {
	return mock.send(queue)
}

// Counts calls to send message methods and returns a new message id. If you pass "error" as queue name it will returns with
// an error and doesn't count this call.
func (mock *publisherMock) send(queue string) (*string, error) {

	if queue == "error" {
		return nil, errors.New("Unable to send message.")
	}

	mock.messageCount++
	messageId := utils.NewId()
	return &messageId, nil
}

// newPublisherMock returns a new mock for a AWS SQS publisher.
func newPublisherMock() sqs.Publisher {
	return &publisherMock{messageCount: 0}
}
