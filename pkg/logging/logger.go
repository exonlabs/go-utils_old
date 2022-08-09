package logging

import (
	"strings"
)

type ILogger interface {
	SetFormatter(format string)
	AddHandler(handler IHandler)
	Info(format string, args ...any)
	Warn(format string, args ...any)
	Debug(format string, args ...any)
	Fatal(format string, args ...any)
}

type Logger struct {
	ILogger
	_handler IHandler
	_name    string
	_fLog    string
}

func GetLogger(name string) ILogger {
	var this Logger
	this.ILogger = &this
	this._name = name
	this._handler = DefaultHandler
	this._handler.SetLogName(name)
	return &this
}

func (this *Logger) SetFormatter(fLog string) {
	var r *strings.Replacer
	var newStr string

	// replace python style by golang style
	r = strings.NewReplacer("%(asctime)s", "%[1]s",
		"%(levelname)s", "%[2]s",
		"%(name)s", "%[3]s",
		"%(message)s", "%[4]s",
	)
	newStr = r.Replace(fLog)

	// build new string if name not added
	if len(this._name) == 0 {
		newStr = r.Replace(fLog)
		newSilceStr := strings.Fields(newStr)
		var b strings.Builder
		for _, s := range newSilceStr {
			if strings.Contains(s, "%[3]s") {
				continue
			}
			b.WriteString(s + " ")
		}
		newStr = strings.TrimSpace(b.String())
	}

	this._fLog = newStr
}

func (this *Logger) AddHandler(handler IHandler) {
	if handler == nil || handler == NewStdoutHandler() {
		this._handler = DefaultHandler
	} else {
		this._handler = handler
	}
	this._handler.SetLogName(this._name)
}

func (this *Logger) Info(format string, args ...any) {
	this._handler.OutPut(this._fLog, "INFO", format, args...)

}

func (this *Logger) Warn(format string, args ...any) {
	this._handler.OutPut(this._fLog, "WARN", format, args...)

}

func (this *Logger) Debug(format string, args ...any) {
	this._handler.OutPut(this._fLog, "DEBUG", format, args...)
}

func (this *Logger) Fatal(format string, args ...any) {
	this._handler.OutPut(this._fLog, "CRITICAL", format, args...)
}
