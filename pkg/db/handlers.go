package db

import (
	"github.com/exonlabs/go-utils/pkg/logging"
)

type IDBHandler interface {
	Session() ISession
	InitializeDB([]IModel)
}

type BaseDBHandler struct {
	IDBHandler
	Options map[string]any
	Backend string
	Logger  *logging.Logger
}

// return session handler object
func (this *BaseDBHandler) Session() ISession {
	panic("NOT_IMPLEMENTED")
}

// create database tables and initialize data
func (this *BaseDBHandler) InitializeDB(models []IModel) {
	// create database structure
	for _, model := range models {
		this.IDBHandler.Session().Query(model).CreateTable()
	}

	// execute migrations
	for _, model := range models {
		model.UpgradeSchema(this.IDBHandler.Session())
	}

	// load initial models data
	for _, model := range models {
		model.InitializeData(this.IDBHandler.Session())
	}
}
