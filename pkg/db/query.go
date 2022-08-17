package db

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type IQuery interface {
	Columns(...string) IQuery
	Filters(string, ...any) IQuery
	FilterBy(string, any) IQuery
	GroupBy(...string) IQuery
	OrderBy(...string) IQuery
	Limit(int32) IQuery
	Offset(int32) IQuery
	DataAdapters(map[string]any) map[string]any
	DataConverters(map[string]any) map[string]any
	Select() []map[string]any
	First() map[string]any
	One() map[string]any
	Count() uint32
	Insert(map[string]any) error
	Update(map[string]any) error
	Delete() error
	CreateTable()
}

type BaseQuery struct {
	IQuery

	DBs   ISession
	Model IModel

	StTablename string
	StColumns   []string
	StFilters   []string
	StExecargs  []any
	StGroupby   []string
	StOrderby   []string
	StLimit     int32
	StOffset    int32
	SqlArgsMap  string
}

// columns: [col1, col2 ...]
func (this *BaseQuery) Columns(columns ...string) IQuery {
	for _, val := range columns {
		this.StColumns = append(this.StColumns, SqlIdentifier(val))
	}
	return this
}

// filters:
func (this *BaseQuery) Filters(filters string, params ...any) IQuery {
	this.StFilters = append(this.StFilters, filters)
	for _, val := range params {
		this.StExecargs = append(this.StExecargs, val)
	}
	return this
}

// filter
func (this *BaseQuery) FilterBy(column string, value any) IQuery {
	if len(this.StFilters) > 0 {
		this.StFilters = append(
			this.StFilters, "AND "+SqlIdentifier(column)+"="+this.SqlArgsMap)
	} else {
		this.StFilters = append(
			this.StFilters, SqlIdentifier(column)+"="+this.SqlArgsMap)
	}
	this.StExecargs = append(this.StExecargs, value)
	return this
}

// groupby: [col1, col2 ...]
func (this *BaseQuery) GroupBy(groupby ...string) IQuery {
	for _, val := range groupby {
		this.StGroupby = append(this.StGroupby, SqlIdentifier(val))
	}
	return this
}

// orderby: [col1 ASC|DESC, col2 ASC|DESC ...]
func (this *BaseQuery) OrderBy(orderby ...string) IQuery {
	for _, val := range orderby {
		v := strings.Split(val, " ")
		v[1] = strings.ToUpper(v[1])
		if v[1] != "ASC" && v[1] != "DESC" {
			panic("invalid sql order type [" + v[1] + "]")
		}
		this.StOrderby = append(
			this.StOrderby, SqlIdentifier(v[0])+" "+v[1])
	}
	return this
}

// limit: integer
func (this *BaseQuery) Limit(limit int32) IQuery {
	this.StLimit = limit
	return this
}

// offset: integer
func (this *BaseQuery) Offset(offset int32) IQuery {
	this.StOffset = offset
	return this
}

func (this *BaseQuery) DataAdapters(data map[string]any) map[string]any {
	model := reflect.Indirect(reflect.ValueOf(this.Model)).
		FieldByName("BaseModel").Interface().(BaseModel)
	for key, fn := range model.DataAdapters {
		if val, ok := data[key]; ok {
			data[key] = fn(val)
		}
	}
	return data
}

func (this *BaseQuery) DataConverters(data map[string]any) map[string]any {
	model := reflect.Indirect(reflect.ValueOf(this.Model)).
		FieldByName("BaseModel").Interface().(BaseModel)
	for key, fn := range model.DataConverters {
		if val, ok := data[key]; ok {
			data[key] = fn(val)
		}
	}
	return data
}

// return all elements matching select query
func (this *BaseQuery) Select() []map[string]any {
	q := "SELECT "
	if len(this.StColumns) > 0 {
		q += strings.Join(this.StColumns, ", ")
	} else {
		q += "*"
	}
	q += " FROM " + SqlIdentifier(this.StTablename)

	if len(this.StFilters) > 0 {
		q += "\nWHERE " + strings.Join(this.StFilters, " ")
	}
	if len(this.StGroupby) > 0 {
		q += "\nGROUP BY " + strings.Join(this.StGroupby, ", ")
	}
	if len(this.StOrderby) > 0 {
		q += "\nORDER BY " + strings.Join(this.StOrderby, ", ")
	}
	if this.StLimit > 0 {
		q += "\nLIMIT " + strconv.Itoa((int)(this.StLimit))
	}
	if this.StOffset > 0 {
		q += "\nOFFSET " + strconv.Itoa((int)(this.StOffset))
	}
	q += ";"
	this.DBs.Connect()
	var data []map[string]any
	for _, val := range this.DBs.FetchAll(q, this.StExecargs...) {
		data = append(data, this.IQuery.DataConverters(val))
	}
	return data
}

// return first element matching select query
func (this *BaseQuery) First() map[string]any {
	return this.Select()[0]
}

// return one element matching select query or nil
//
//	there must be only one element if exists
func (this *BaseQuery) One() map[string]any {
	res := this.Limit(2).Select()
	if res != nil {
		if len(res) > 2 {
			panic("multiple entries found")
		}
		return res[0]
	}
	return nil
}

func (this *BaseQuery) Count() uint32 {
	q := "SELECT count(*) as count FROM " + SqlIdentifier(this.StTablename)
	if len(this.StFilters) > 0 {
		q += "\nWHERE " + strings.Join(this.StFilters, " ")
	}
	if len(this.StGroupby) > 0 {
		q += "\nGROUP BY " + strings.Join(this.StGroupby, ", ")
	}
	q += ";"
	this.DBs.Connect()
	data := this.DBs.FetchOne(q)
	return uint32(data["count"].(int64))
}

func (this *BaseQuery) Insert(data map[string]any) error {
	var columns []string
	var params []any

	if _, ok := data["guid"]; !ok {
		columns = append(columns, "guid")
		params = append(params, this.generateGuid())
	}

	data = this.IQuery.DataAdapters(data)

	for k, v := range data {
		columns = append(columns, SqlIdentifier(k))
		params = append(params, v)
	}

	q := "INSERT INTO " + this.StTablename
	q += fmt.Sprintf("\n(%s)", strings.Join(columns, ", "))
	q += "\nVALUES"
	q += fmt.Sprintf("\n(%s", strings.Repeat(this.SqlArgsMap+" ,", len(columns)))
	q = strings.TrimSuffix(q, ",")
	q += ")"
	q += ";"

	this.DBs.Connect()
	return this.DBs.Execute(q, params...)
}

func (this *BaseQuery) Update(data map[string]any) error {
	var columns []string
	var params []any

	if _, ok := data["guid"]; ok {
		delete(data, "guid")
	}

	data = this.IQuery.DataAdapters(data)

	for k, v := range data {
		columns = append(columns, SqlIdentifier(k))
		params = append(params, v)
	}

	q := "UPDATE " + this.StTablename
	q += fmt.Sprintf("\nSET %s", strings.Join(columns, "= "+this.SqlArgsMap+", "))
	q += "= " + this.SqlArgsMap
	if len(this.StFilters) > 0 {
		q += "\nWHERE " + strings.Join(this.StFilters, " ")
		params = append(params, this.StExecargs...)
	}

	this.DBs.Connect()
	return this.DBs.Execute(q, params...)
}

func (this *BaseQuery) Delete() error {
	q := "DELETE FROM " + this.StTablename
	if len(this.StFilters) > 0 {
		q += "\nWHERE " + strings.Join(this.StFilters, " ")
	}

	this.DBs.Connect()
	return this.DBs.Execute(q, this.StExecargs...)
}

func (this *BaseQuery) CreateTable() {
	panic("NOT_IMPLEMENTED")
}

func (this *BaseQuery) generateGuid() string {
	return hex.EncodeToString(uuid.NewV5(uuid.NewV1(), uuid.NewV4().String()).Bytes())
}
