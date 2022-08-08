package db

import (
	"strconv"
	"strings"
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
	// First() *map[string]any
	// One() *map[string]any
	// Count() uint32
	Insert(map[string]any) error
	// Update(, map[string]any) uint32
	// Delete() uint32
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
	return this.DBs.FetchAll(q)
}

func (this *BaseQuery) Insert(data map[string]any) error {
	var values []any
	q := "INSERT INTO " + this.Model.TableName() + "("
	for k, v := range data {
		q += k
		q += ", "
		values = append(values, v)
	}
	q = strings.TrimSpace(q)
	q = strings.TrimSuffix(q, ",")
	q += ") VALUES ("
	q += strings.Repeat("? ,", len(data))
	q = strings.TrimSuffix(q, ",")
	q += ")"
	this.DBs.Connect()
	return this.DBs.Execute(q, values...)
}

func (this *BaseQuery) CreateTable() {
	panic("NOT_IMPLEMENTED")
}
