package handlers

import (
	"os"

	"github.com/exonlabs/go-utils/logging"
)

type FileHandler struct {
	*BaseHandler
	Filepath string
}

func NewFileHandler(filepath string) *FileHandler {
	return &FileHandler{
		BaseHandler: NewBaseHandler(),
		Filepath:    filepath,
	}
}

func NewFileLogger(name string, filepath string) *logging.Logger {
	logger := logging.NewLogger(name)
	logger.AddHandler(NewFileHandler(filepath))
	return logger
}

func (hnd *FileHandler) HandleRecord(
	logger *logging.Logger, record *logging.Record) {
	if msg := hnd.ProcessRecord(logger, record); msg != "" {
		hnd.EmitMessage(msg)
	}
}

func (hnd *FileHandler) EmitMessage(msg string) {
	file, err := os.OpenFile(
		hnd.Filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	_, err = file.Write([]byte(msg + "\n"))
	if err != nil {
		panic(err)
	}
}
