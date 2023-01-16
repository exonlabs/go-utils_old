package handlers

import (
	"strings"

	"github.com/exonlabs/go-utils/logging"
)

type BaseHandler struct {
	Level     logging.LogLevel
	Formatter string
}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{
		Level: logging.LEVEL_NOTSET,
	}
}

func (hnd *BaseHandler) SetFormatter(frmt string) {
	hnd.Formatter = frmt
}

func (hnd *BaseHandler) ProcessRecord(
	logger *logging.Logger, record *logging.Record) string {

	filterlvl := hnd.Level
	if filterlvl == logging.LEVEL_NOTSET {
		filterlvl = logger.Level
	}
	if filterlvl == logging.LEVEL_NOTSET || record.Level < filterlvl {
		return ""
	}

	r := strings.NewReplacer(
		"%(asctime)s", record.Timestamp.Format(
			"2006-01-02 15:04:05.000000"),
		"%(name)s", record.Caller,
		"%(levelname)s", record.LevelName,
		"%(message)s", strings.TrimSpace(record.Message),
	)
	if hnd.Formatter == "" {
		return r.Replace(logger.Formatter)
	} else {
		return r.Replace(hnd.Formatter)
	}
}
