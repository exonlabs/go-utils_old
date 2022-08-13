package logging

import (
	"fmt"
	"time"
)

var (
	FATAL   uint8 = 50
	ERROR   uint8 = 40
	WARNING uint8 = 30
	INFO    uint8 = 20
	DEBUG   uint8 = 10
	NOTSET  uint8 = 0
)

type Record struct {
	Timestamp time.Time
	Name      string
	Level     uint8
	Message   string
}

type Logger struct {
	Name     string
	Level    uint8
	Handlers []IHandler
}

func NewLogger() *Logger {
	var this Logger
	this.Level = INFO
	this.Handlers = []IHandler{NewStdoutHandler()}
	return &this
}

func (this *Logger) AddHandler(handler IHandler) {
	this.Handlers = append(this.Handlers, handler)
}

func (this *Logger) ClearHandlers() {
	this.Handlers = []IHandler{}
}

func (this *Logger) SetFormatter(fmtstr string) {
	for _, hnd := range this.Handlers {
		hnd.SetFormatter(fmtstr)
	}
}

func (this *Logger) Log(level uint8, msg string, args ...any) {
	record := Record{
		Timestamp: time.Now().Local(),
		Name:      this.Name,
		Level:     level,
		Message:   fmt.Sprintf(msg, args...),
	}
	for _, hnd := range this.Handlers {
		hnd.Handle(this, &record)
	}
}

func (this *Logger) Debug(msg string, args ...any) {
	this.Log(DEBUG, msg, args...)
}

func (this *Logger) Info(msg string, args ...any) {
	this.Log(INFO, msg, args...)
}

func (this *Logger) Warn(msg string, args ...any) {
	this.Log(WARNING, msg, args...)
}

func (this *Logger) Error(msg string, args ...any) {
	this.Log(ERROR, msg, args...)
}

func (this *Logger) Fatal(msg string, args ...any) {
	this.Log(FATAL, msg, args...)
}
