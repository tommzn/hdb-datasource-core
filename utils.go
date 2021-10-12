package core

import (
	"errors"
	"strings"

	config "github.com/tommzn/go-config"
)

// asError returns a single error with all passed or nil if passed slice is empty.
func asError(errorList []error) error {
	if len(errorList) > 0 {
		errorMessages := []string{}
		for _, err := range errorList {
			errorMessages = append(errorMessages, err.Error())
		}
		return errors.New(strings.Join(errorMessages, "\n"))
	}
	return nil
}

// archiveQueueFromConfig try to get archive queue name from passed config
// and will return with a default value if nothing can be obtained.
func archiveQueueFromConfig(conf config.Config) string {
	archiveQueue := conf.Get("hdb.archive", config.AsStringPtr("de.tsl.hdb.archive"))
	return *archiveQueue
}
