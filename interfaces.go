package core

// DataSource retrieves data from a specific source.
type DataSource interface {

	// Fetch will retrieve new data and returns it as an event from central event lib.
	// For more details about events see https://github.com/tommzn/hdb-events-go
	Fetch() (interface{}, error)
}
