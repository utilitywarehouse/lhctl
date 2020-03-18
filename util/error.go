package util

import (
	"log"
	"strings"
)

type ErrorHandler struct {
}

type ErrorHandlerInterface interface {
	ExitOnError(error, ...string)
}

func (eh *ErrorHandler) ExitOnError(err error, msg ...string) {
	errMsg := "error"
	if len(msg) > 0 {
		errMsg = strings.Join(msg, ", ")
	}

	log.Fatalf("%s: %v", errMsg, err)
}
