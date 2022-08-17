package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/exonlabs/go-utils/pkg/db"
)

type Session struct {
	db.BaseSession
}

func NewSession(options map[string]any) *Session {
	var this Session
	this.ISession = &this
	this.Options = options
	this.Is_Connected = false
	return &this
}

func (this *Session) Query(model db.IModel) db.IQuery {
	return NewQuery(this, model)
}

func (this *Session) Connect() {
	if val, ok := this.Options["database"]; ok {
		db, err := sql.Open("sqlite3", val.(string))
		if err != nil {
			panic(err)
		}
		this.SqlDB = db
	} else {
		// todo log message v
		panic("invalid database configuration")
	}

	this.Is_Connected = true
}
