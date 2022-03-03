package core

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/golang/protobuf/proto"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
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
	archiveQueue := conf.Get("hdb.archive", config.AsStringPtr("de-tsl-hdb-archive"))
	return *archiveQueue
}

func logEvent(event proto.Message, logger log.Logger) {

	logger.Debugf("EVent: %+v", event)
	protoContent, _ := proto.Marshal(event)
	logger.Debugf("Proto: %s", string(protoContent))
	logger.Debugf("Base64: %s", base64.StdEncoding.EncodeToString(protoContent))
}
