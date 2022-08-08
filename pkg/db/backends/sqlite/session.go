package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/exonlabs/go-utils/pkg/db"
)

type Session struct {
	db.BaseSession
}

func NewSession(options map[string]any, logger any, debug uint8) db.ISession {
	var this Session
	this.ISession = &this
	this.Backend = "sqlite"
	this.Options = options
	this.Debug = debug
	this.Is_Connected = false

	// TODO: implement real logger
	// this.Logger = logger
	// switch {
	//     case this.Debug >= 5:
	//         this.Logger.SetLevel("DEBUG")
	//     case this.Debug >= 4:
	//         this.Logger.SetLevel("INFO")
	//     default:
	//         this.Logger.SetLevel("ERROR")
	// }

	return &this
}

func (this *Session) Query(model db.IModel) *db.BaseQuery {
	return &db.BaseQuery{
		DBs:   this,
		Model: model,
	}
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
