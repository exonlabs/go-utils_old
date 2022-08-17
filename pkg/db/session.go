package db

import (
	"database/sql"

	"github.com/exonlabs/go-utils/pkg/logging"
)

type ISession interface {
	Query(IModel) IQuery
	Connect()
	Close()
	IsConnected() bool
	Execute(string, ...any) error
	FetchOne(string, ...any) map[string]any
	FetchAll(string, ...any) []map[string]any
	RowsAffected() int64
	LastInsertId() int64
	Begin()
	Commit()
	RollBack()
}

type BaseSession struct {
	ISession
	Options      map[string]any
	Logger       *logging.Logger
	SqlDB        *sql.DB
	SqlTX        *sql.Tx
	Is_Connected bool
}

func (this *BaseSession) LogSql(query string, params ...any) {
	if this.Logger != nil && this.Logger.Level == logging.DEBUG {
		this.Logger.Debug("SQL:\n---\n" + query + "\nPARAMS: %v\n---", params...)
	}
}

func (this *BaseSession) Query(model IModel) IQuery {
	panic("NOT_IMPLEMENTED")
}

func (this *BaseSession) IsConnected() bool {
	return this.Is_Connected
}

func (this *BaseSession) Connect() {
	panic("NOT_IMPLEMENTED")
}

func (this *BaseSession) Close() {
	if this.Is_Connected {
		if err := this.SqlDB.Close(); err != nil {
			panic(err)
		}
		this.Is_Connected = false
	}
}

func (this *BaseSession) Execute(sql string, params ...any) error {
	if !this.Is_Connected {
		panic("not connected")
	}

	if this.SqlTX != nil {
		_, err := this.SqlTX.Exec(sql, params...)
		if err != nil {
			return err
		}
	} else {
		_, err := this.SqlDB.Exec(sql, params...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *BaseSession) FetchOne(sql string, params ...any) map[string]any {
	if !this.Is_Connected {
		panic("not connected")
	}

	return this.FetchAll(sql, params...)[0]
}

func (this *BaseSession) FetchAll(sql string, params ...any) []map[string]any {
	if !this.Is_Connected {
		panic("not connected")
	}

	rows, err := this.SqlDB.Query(sql, params...)
	if err != nil {
		panic(err)
	}

	// get columns name
	cols, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	// create slice of map to fill data
	var data []map[string]any

	for rows.Next() {
		// create a slice of any's to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]any, len(cols))
		columnPointers := make([]any, len(cols))
		for k := range columns {
			columnPointers[k] = &columns[k]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			panic(err)
		}
		// create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		d := make(map[string]any, len(cols))
		for k, colName := range cols {
			val := columnPointers[k].(*any)
			d[colName] = *val
		}
		data = append(data, d)
	}
	rows.Close()

	return data
}

func (this *BaseSession) RowsAffected() int64 {
	panic("NOT_IMPLEMENTED")
}

func (this *BaseSession) LastInsertId() int64 {
	panic("NOT_IMPLEMENTED")
}

func (this *BaseSession) Begin() {
	this.ISession.Connect()
	tx, err := this.SqlDB.Begin()
	if err != nil {
		panic(err)
	}
	this.SqlTX = tx
}

func (this *BaseSession) Commit() {
	if !this.Is_Connected {
		panic("not connected")
	}
	if err := this.SqlTX.Commit(); err != nil {
		panic(err)
	}
	this.SqlTX = nil
}

func (this *BaseSession) RollBack() {
	if !this.Is_Connected {
		panic("not connected")
	}
	if err := this.SqlTX.Rollback(); err != nil {
		panic(err)
	}
	this.SqlTX = nil
}
