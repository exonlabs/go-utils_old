package sqlite

import (
	"github.com/exonlabs/go-utils/pkg/db"
)

type DBHandler struct {
	db.BaseDBHandler
}

func NewDBHandler(options map[string]any) *DBHandler {
	var this DBHandler
	this.IDBHandler = &this
	this.Options = options
	this.Backend = "sqllite"
	return &this
}

func (this *DBHandler) Session() db.ISession {
	sess := NewSession(this.Options)
	sess.Logger = this.Logger
	return sess
}
