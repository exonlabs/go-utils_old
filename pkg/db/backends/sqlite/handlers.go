package sqlite

import (
	"github.com/exonlabs/go-utils/pkg/db"
)

type DBHandler struct {
	db.BaseDBHandler
}

func NewDBHandler(options map[string]any) db.IDBHandler {
	var this DBHandler
	this.IDBHandler = &this
	this.Logger = nil
	this.Options = options
	this.Backend = "sqllite"
	return &this
}

func (this *DBHandler) CreateSession() db.ISession {
	return NewSession(this.Options, this.Logger, 1)
}

// create database tables and initialize data
func (this *DBHandler) InitDatabase(models []db.IModel) bool {
	for _, model := range models {
		q := NewQuery(this.IDBHandler.CreateSession(), model)
		q.CreateTable()
	}
	for _, model := range models {
		model.InitializeData(this.IDBHandler.CreateSession())
	}
	return true
}
