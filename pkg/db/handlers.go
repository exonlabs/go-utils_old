package db

import (
	"fmt"
)


type IDBHandler interface {
	CreateSession() *ISession
	InitDatabase([]Model) bool
}

type BaseDBHandler struct {
	IDBHandler

	Options map[string]any
	Logger any
	Backend string
}


// return session handler object
func (this *BaseDBHandler) CreateSession() *ISession {
	panic("NOT_IMPLEMENTED")
}

// create database tables and initialize data
func (this *BaseDBHandler) InitDatabase(models []Model) bool {
	return true
}
