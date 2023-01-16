package handlers

import (
	"fmt"

	"github.com/exonlabs/go-utils/logging"
)

type StdoutHandler struct {
	*BaseHandler
}

func NewStdoutHandler() *StdoutHandler {
	return &StdoutHandler{
		BaseHandler: NewBaseHandler(),
	}
}

func NewStdoutLogger(name string) *logging.Logger {
	logger := logging.NewLogger(name)
	logger.AddHandler(NewStdoutHandler())
	return logger
}

func (hnd *StdoutHandler) HandleRecord(
	logger *logging.Logger, record *logging.Record) {
	if msg := hnd.ProcessRecord(logger, record); msg != "" {
		hnd.EmitMessage(msg)
	}
}

func (hnd *StdoutHandler) EmitMessage(msg string) {
	fmt.Print(msg + "\n")
}
