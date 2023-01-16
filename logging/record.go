package logging

import (
	"fmt"
	"time"
)

type Record struct {
	Timestamp time.Time
	Caller    string
	Level     LogLevel
	LevelName string
	Message   string
}

// create new logging record
func NewRecord(
	caller string, level LogLevel, msg string, args ...any) *Record {

	var lvlname string
	switch level {
	case LEVEL_FATAL:
		lvlname = "FATAL"
	case LEVEL_ERROR:
		lvlname = "ERROR"
	case LEVEL_WARN:
		lvlname = "WARN "
	case LEVEL_INFO:
		lvlname = "INFO "
	default:
		lvlname = "DEBUG"
	}

	return &Record{
		Timestamp: time.Now().Local(),
		Caller:    caller,
		Level:     level,
		LevelName: lvlname,
		Message:   fmt.Sprintf(msg, args...),
	}
}
