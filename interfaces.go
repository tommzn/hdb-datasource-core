package core

import "github.com/golang/protobuf/proto"

// DataSource retrieves data from a specific source.
type DataSource interface {

	// Fetch will retrieve new data and returns it as an event from central event lib.
	// For more details about events see https://github.com/tommzn/hdb-events-go
	Fetch() (proto.Message, error)
}

// A Collector calls fetch method of a datasource and process the returned event.
type Collector interface {

	// Run executes collectors data processing logic.
	Run() error
}
