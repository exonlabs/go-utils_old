package logging

import (
	"fmt"
	"os"
	"strings"
)

type IHandler interface {
	SetFormatter(string)
	Handle(*Logger, *Record)
	Emit(string)
}

type BaseHandler struct {
	IHandler
	Level     uint8
	Formatter string
}

func (this *BaseHandler) SetFormatter(fmtstr string) {
	this.Formatter = fmtstr
}

func (this *BaseHandler) Handle(logger *Logger, record *Record) {
	filter_lvl := logger.Level
	if this.Level > NOTSET {
		filter_lvl = this.Level
	}
	if filter_lvl > NOTSET && record.Level >= filter_lvl {
		var lvlname string
		switch record.Level {
		case FATAL:
			lvlname = "FATAL"
		case ERROR:
			lvlname = "ERROR"
		case WARNING:
			lvlname = "WARN "
		case INFO:
			lvlname = "INFO "
		case DEBUG:
			lvlname = "DEBUG"
		}

		var r *strings.Replacer
		r = strings.NewReplacer(
			"%(asctime)s", record.Timestamp.Format("2006-01-02 15:04:05.000000"),
			"%(name)s", record.Name,
			"%(levelname)s", lvlname,
			"%(message)s", strings.TrimSpace(record.Message),
		)
		logline := r.Replace(this.Formatter)
		this.IHandler.Emit(logline + "\n")
	}
}

func (this *BaseHandler) Emit(logline string) {
	panic("NOT_IMPLEMENTED")
}

type StdoutHandler struct {
	BaseHandler
}

func NewStdoutHandler() *StdoutHandler {
	var this StdoutHandler
	this.IHandler = &this
	this.Level = NOTSET
	this.Formatter = "%(asctime)s %(levelname)s %(message)s"
	return &this
}

func (this *StdoutHandler) Emit(logline string) {
	fmt.Print(logline)
}

type FileHandler struct {
	BaseHandler
	Filepath string
}

func NewFileHandler(filepath string) *FileHandler {
	var this FileHandler
	this.IHandler = &this
	this.Level = NOTSET
	this.Formatter = "%(asctime)s %(levelname)s %(message)s"
	this.Filepath = filepath
	return &this
}

func (this *FileHandler) Emit(logline string) {
	file, err := os.OpenFile(
		this.Filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	_, err = file.Write([]byte(logline))
	if err != nil {
		panic(err)
	}
}
