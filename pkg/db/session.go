package db

import (
	"database/sql"
	"fmt"
	"log"
)

type ISession interface {
	LogSql(string, ...any)
	Query(IModel) *BaseQuery
	IsConnected() bool
	Connect()
	Close()
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

	SqlDB        *sql.DB
	sqlTX        *sql.Tx
	Backend      string
	Options      map[string]any
	Logger       any
	Debug        uint8
	Is_Connected bool
}

func (this *BaseSession) LogSql(query string, params ...any) {
	// TODO: use real logger
	fmt.Printf("SQL:\n---\n%v\n---", query)
	if len(params) > 0 {
		fmt.Printf("PARAMS: %v\n", params)
	}
}

func (this *BaseSession) Query(model IModel) *BaseQuery {
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

	if this.sqlTX != nil {
		_, err := this.sqlTX.Exec(sql, params...)
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
		log.Println("All Rows Error:", err)
	}

	// get columns name
	cols, err := rows.Columns()
	if err != nil {
		log.Println("All Cols Error:", err)
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
			log.Println("All Scan Error:", err)
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
	this.sqlTX = tx
}

func (this *BaseSession) Commit() {
	if !this.Is_Connected {
		panic("not connected")
	}
	if err := this.sqlTX.Commit(); err != nil {
		panic(err)
	}
	this.sqlTX = nil
}

func (this *BaseSession) RollBack() {
	if !this.Is_Connected {
		panic("not connected")
	}
	if err := this.sqlTX.Rollback(); err != nil {
		panic(err)
	}
	this.sqlTX = nil
}
