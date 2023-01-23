package db

import (
	"github.com/exonlabs/go-utils/logging"
)

type Handler struct {
	Engine  IEngine
	Options KwArgs
	Logger  *logging.Logger
	Debug   logging.LogLevel
}

func NewHandler(engine IEngine, options KwArgs) *Handler {
	// set default options
	if _, ok := options["connect_timeout"].(uint); !ok {
		options["connect_timeout"] = 30
	}
	if _, ok := options["retries"].(uint); !ok {
		options["retries"] = 10
	}
	if _, ok := options["retry_delay"].(float32); !ok {
		options["retry_delay"] = 0.5
	}

	return &Handler{
		Engine:  engine,
		Options: options,
		Debug:   logging.LEVEL_ERROR,
	}
}

// create new session from handler
func (dbh *Handler) Session() *Session {
	return NewSession(dbh)
}

// create database tables and initialize data
func (dbh *Handler) InitDatabase(tables map[TableName]IModel) error {
	dbs := dbh.Session()
	defer dbs.Close()

	// create database schema in session transaction
	if err := dbs.Begin(); err != nil {
		return err
	}
	// 1st create all tables if not exist
	for tblname, model := range tables {
		schema, err := dbh.Engine.GenTableSchema(
			tblname, model.GetTableMeta())
		if err != nil {
			return err
		}
		for _, stmt := range schema {
			if err := dbs.Execute(stmt); err != nil {
				return err
			}
		}
	}
	// 2nd upgrade previous tables
	for tblname, model := range tables {
		if m, ok := model.(IModelUpgradeTableSchema); ok {
			if err := m.UpgradeTableSchema(dbs, tblname); err != nil {
				return err
			}
		}
	}
	if err := dbs.Commit(); err != nil {
		return err
	}

	// 3rd initialize models data
	for tblname, model := range tables {
		if m, ok := model.(IModelInitializeData); ok {
			if err := m.InitializeData(dbs, tblname); err != nil {
				return err
			}
		}
	}

	return nil
}
