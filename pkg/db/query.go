package db

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type IQuery interface {
	Columns(columns ...string) IQuery
	Filters(filters string, params ...any) IQuery
	FilterBy(column string, value any) IQuery
	GroupBy(groupby ...string) IQuery
	OrderBy(orderby ...string) IQuery
	Limit(limit int32) IQuery
	Offset(offset int32) IQuery
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

	_tablename string
	_columns   []string
	_filters   []string
	_execargs  []any
	_groupby   []string
	_orderby   []string
	_limit     int32
	_offset    int32
}

// columns: [col1, col2 ...]
func (this *BaseQuery) Columns(columns ...string) IQuery {
	for _, val := range columns {
		this._columns = append(this._columns, SqlIdentifier(val))
	}
	return this
}

// filters:
func (this *BaseQuery) Filters(filters string, params ...any) IQuery {
	this._filters = append(this._filters, filters)
	for _, val := range params {
		this._execargs = append(this._execargs, val)
	}
	return this
}

// filter
func (this *BaseQuery) FilterBy(column string, value any) IQuery {
	if len(this._filters) > 0 {
		this._filters = append(
			this._filters, "AND "+SqlIdentifier(column)+"=?")
	} else {
		this._filters = append(
			this._filters, SqlIdentifier(column)+"=?")
	}
	this._execargs = append(this._execargs, value)
	return this
}

// groupby: [col1, col2 ...]
func (this *BaseQuery) GroupBy(groupby ...string) IQuery {
	for _, val := range groupby {
		this._groupby = append(this._groupby, SqlIdentifier(val))
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
		this._orderby = append(
			this._orderby, SqlIdentifier(v[0])+" "+v[1])
	}
	return this
}

// limit: integer
func (this *BaseQuery) Limit(limit int32) IQuery {
	this._limit = limit
	return this
}

// offset: integer
func (this *BaseQuery) Offset(offset int32) IQuery {
	this._offset = offset
	return this
}

// return all elements matching select query
func (this *BaseQuery) Select() []map[string]any {
	q := "SELECT "
	if len(this._columns) > 0 {
		q += strings.Join(this._columns, ", ")
	} else {
		q += "*"
	}
	q += " FROM " + SqlIdentifier(this.Model.TableName())

	if len(this._filters) > 0 {
		q += "\nWHERE " + strings.Join(this._filters, " ")
	}
	if len(this._groupby) > 0 {
		q += "\nGROUP BY " + strings.Join(this._groupby, ", ")
	}
	if len(this._orderby) > 0 {
		q += "\nORDER BY " + strings.Join(this._orderby, ", ")
	}
	if this._limit > 0 {
		q += "\nLIMIT " + strconv.Itoa((int)(this._limit))
	}
	if this._offset > 0 {
		q += "\nOFFSET " + strconv.Itoa((int)(this._offset))
	}
	q += ";"
	this.DBs.Connect()
	return this.DBs.FetchAll(q, this._execargs...)
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
	q := "SELECT count(*) as count FROM " + SqlIdentifier(this.Model.TableName())
	if len(this._filters) > 0 {
		q += "\nWHERE " + strings.Join(this._filters, " ")
	}
	if len(this._groupby) > 0 {
		q += "\nGROUP BY " + strings.Join(this._groupby, ", ")
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

	data = this.Model.DataAdapters(data)

	for k, v := range data {
		columns = append(columns, SqlIdentifier(k))
		params = append(params, v)
	}

	q := "INSERT INTO " + this.Model.TableName()
	q += fmt.Sprintf("\n(%s)", strings.Join(columns, ", "))
	q += "\nVALUES"
	q += fmt.Sprintf("\n(%s", strings.Repeat("? ,", len(columns)))
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

	data = this.Model.DataAdapters(data)

	for k, v := range data {
		columns = append(columns, SqlIdentifier(k))
		params = append(params, v)
	}

	q := "UPDATE " + this.Model.TableName()
	q += fmt.Sprintf("\nSET %s", strings.Join(columns, "= ?, "))
	q += "= ?"
	if len(this._filters) > 0 {
		q += "\nWHERE " + strings.Join(this._filters, " ")
		params = append(params, this._execargs...)
	}

	this.DBs.Connect()
	return this.DBs.Execute(q, params...)
}

func (this *BaseQuery) Delete() error {
	q := "DELETE FROM " + this.Model.TableName()
	if len(this._filters) > 0 {
		q += "\nWHERE " + strings.Join(this._filters, " ")
	}

	this.DBs.Connect()
	return this.DBs.Execute(q, this._execargs...)
}

func (this *BaseQuery) CreateTable() {
	panic("NOT_IMPLEMENTED")
}

func (this *BaseQuery) generateGuid() string {
	return hex.EncodeToString(uuid.NewV5(uuid.NewV1(), uuid.NewV4().String()).Bytes())
}
