package core

import (
	"context"
	"errors"
	"sync"
	"time"

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
	}
	return &events.Dummy{Id: utils.NewId()}, nil
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
func (mock *publisherMock) Send(message proto.Message) error {

	if mock.shouldFail {
		return errors.New("Unable to send message.")
	}

	mock.messageCount += 2
	return nil
}

// newPublisherMock returns a new mock for a AWS SQS publisher.
func newPublisherMock() Publisher {
	return &publisherMock{shouldFail: false, messageCount: 0}
}

// s3EventProcessorMock can be used for testing, processing will aalways succeed.
type s3EventProcessorMock struct {
}

// ProcessEvent will return always no message and no error.
func (mock *s3EventProcessorMock) ProcessEvent(entity awsevents.S3Entity, content []byte) (proto.Message, error) {
	return nil, nil
}

// newS3EventProcessorMock returns a newS3 event processor mock for testing.
func newS3EventProcessorMock() S3EventProcessor {
	return &s3EventProcessorMock{}
}

// collectorMock is used to test continuous collector.
type collectorMock struct {
	counter     int
	hasFinished bool
}

// Run will increade a counter every 100 milliseconds and stops execution
// if context has been canceled.
func (mock *collectorMock) Run(ctx context.Context) error {

	stopChan := make(chan bool, 1)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		mock.counter++
		time.Sleep(100 * time.Millisecond)
		if len(stopChan) > 0 {
			return
		}
	}()
	<-ctx.Done()
	stopChan <- true
	wg.Wait()
	mock.hasFinished = true
	return nil
}

// newCollectorMock returns a new collector mock for testing.
func newCollectorMock() Collector {
	return &collectorMock{
		counter:     0,
		hasFinished: false,
	}
}
