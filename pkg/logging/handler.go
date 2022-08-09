package logging

import (
	"fmt"
	"os"
)

var DefaultHandler = NewStdoutHandler()

type IHandler interface {
	OutPut(fLog string, level string, format string, args ...any)
	SetLogName(name string)
	GetLogName() string
}

type StdoutHandler struct {
	IHandler
	_logName string
}

func NewStdoutHandler() IHandler {
	var this StdoutHandler
	this.IHandler = &this
	return &this
}

func (this *StdoutHandler) OutPut(fLog string, level string, format string, args ...any) {
	fmt.Print(formater(fLog, this.GetLogName(), level, format, args...))
}

func (this *StdoutHandler) SetLogName(name string) {
	this._logName = name
}

func (this *StdoutHandler) GetLogName() string {
	return this._logName
}

type FileHandler struct {
	IHandler
	_path    string
	_logName string
}

func NewFileHandler(path string) IHandler {
	var this FileHandler
	this.IHandler = &this
	this._path = path
	return &this
}

func (this *FileHandler) OutPut(fLog string, level string, format string, args ...any) {
	file, err := os.OpenFile(this._path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = file.Write([]byte(formater(fLog, this.GetLogName(), level, format, args...)))
	if err != nil {
		panic(err)
	}
}

func (this *FileHandler) SetLogName(name string) {
	this._logName = name
}

func (this *FileHandler) GetLogName() string {
	return this._logName
}
