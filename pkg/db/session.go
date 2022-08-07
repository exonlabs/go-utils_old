package db

import (
	"fmt"
)


type ISession interface {
	LogSql(string, ...any)
	Query() *Query
	IsConnected() bool
	Connect()
	Close()
	Execute(string, ...any)
	FetchOne() *map[string]any
	FetchAll() *[]map[string]any
	RowCount() int32
	LastRowId() int32
	Commit()
	RollBack()
}

type BaseSession struct {
	ISession

	Backend string
	Options map[string]any
	Logger any
	Debug uint8
	IsConnected bool
}


func (this *BaseSession) LogSql(query string, params ...any) {
	// TODO: use real logger
	fmt.Printf("SQL:\n---\n%v\n---", query)
	if len(params) > 0 {
		fmt.Printf("PARAMS: %v\n", params)
	}
}

func (this *BaseSession) Query() *Query {
	panic("NOT_IMPLEMENTED")
}

func (this *BaseSession) IsConnected() bool {
	panic("NOT_IMPLEMENTED")
}

func (this *BaseSession) Connect() {
	panic("NOT_IMPLEMENTED")
}

func (this *BaseSession) Close() {
	panic("NOT_IMPLEMENTED")
}

func (this *BaseSession) Execute(sql string, params ...any) {
	panic("NOT_IMPLEMENTED")
}

func (this *BaseSession) FetchOne() *map[string]any {
	panic("NOT_IMPLEMENTED")
}

func (this *BaseSession) FetchAll() *[]map[string]any {
	panic("NOT_IMPLEMENTED")
}

func (this *BaseSession) RowCount() int32 {
	panic("NOT_IMPLEMENTED")
}

func (this *BaseSession) LastRowId() int32 {
	panic("NOT_IMPLEMENTED")
}

func (this *BaseSession) Commit() {
	panic("NOT_IMPLEMENTED")
}

func (this *BaseSession) RollBack() {
	panic("NOT_IMPLEMENTED")
}
