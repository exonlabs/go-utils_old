package logging

type LogLevel = uint8

const (
	// logging levels
	LEVEL_NOTSET LogLevel = 0
	LEVEL_DEBUG  LogLevel = 10
	LEVEL_INFO   LogLevel = 20
	LEVEL_WARN   LogLevel = 30
	LEVEL_ERROR  LogLevel = 40
	LEVEL_FATAL  LogLevel = 50
)

type IHandler interface {
	SetFormatter(string)
	HandleRecord(*Logger, *Record) error
	EmitMessage(string) error
}
