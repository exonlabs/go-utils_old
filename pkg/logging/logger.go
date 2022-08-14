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
	Parent   *Logger
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

func (this *Logger) Log(record *Record) {
	for _, hnd := range this.Handlers {
		hnd.Handle(this, record)
	}

	// parent log
	if this.Parent != nil {
		this.Parent.Log(record)
	}
}

func (this *Logger) Debug(msg string, args ...any) {
	record := Record{
		Timestamp: time.Now().Local(),
		Name:      this.Name,
		Level:     DEBUG,
		Message:   fmt.Sprintf(msg, args...),
	}
	this.Log(&record)
}

func (this *Logger) Info(msg string, args ...any) {
	record := Record{
		Timestamp: time.Now().Local(),
		Name:      this.Name,
		Level:     INFO,
		Message:   fmt.Sprintf(msg, args...),
	}
	this.Log(&record)
}

func (this *Logger) Warn(msg string, args ...any) {
	record := Record{
		Timestamp: time.Now().Local(),
		Name:      this.Name,
		Level:     WARNING,
		Message:   fmt.Sprintf(msg, args...),
	}
	this.Log(&record)
}

func (this *Logger) Error(msg string, args ...any) {
	record := Record{
		Timestamp: time.Now().Local(),
		Name:      this.Name,
		Level:     ERROR,
		Message:   fmt.Sprintf(msg, args...),
	}
	this.Log(&record)
}

func (this *Logger) Fatal(msg string, args ...any) {
	record := Record{
		Timestamp: time.Now().Local(),
		Name:      this.Name,
		Level:     FATAL,
		Message:   fmt.Sprintf(msg, args...),
	}
	this.Log(&record)
}
