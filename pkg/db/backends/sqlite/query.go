package sqlite

import (
	"strconv"
	"strings"
)


type Query struct {
	dbs Session
	_tablename string
	_columns []string
	_filters []string
	_execargs []any
	_groupby []string
	_orderby []string
	_limit int32
	_offset int32
}

func NewQuery(dbs Session, tablename string) *Query {
	return &Query{
		dbs: dbs,
		_tablename: tablename,
		_columns: []string{},
		_filters: []string{},
		_execargs: []any{},
		_groupby: []string{},
		_orderby: []string{},
		_limit: -1,
		_offset: -1,
	}
}

// columns: [col1, col2 ...]
func (this *Query) Columns(columns ...string) *Query {
	for _, val := range columns {
		this._columns = append(this._columns, SqlIdentifier(val))
	}
	return this
}

// filters:
func (this *Query) Filters(filters string, params ...any) *Query {
	this._filters = append(this._filters, filters)
	for _, val := range params {
		this._execargs = append(this._execargs, val)
	}
	return this
}

// filter
func (this *Query) FilterBy(column string, value any) *Query {
	if len(this._filters) > 0 {
		this._filters = append(
			this._filters, "AND " + SqlIdentifier(column) + "=?")
	} else {
		this._filters = append(
			this._filters, SqlIdentifier(column) + "=?")
	}
	this._execargs = append(this._execargs, value)
	return this
}

// groupby: [col1, col2 ...]
func (this *Query) GroupBy(groupby ...string) *Query {
	for _, val := range groupby {
		this._groupby = append(this._groupby, SqlIdentifier(val))
	}
	return this
}

// orderby: [col1 ASC|DESC, col2 ASC|DESC ...]
func (this *Query) OrderBy(orderby ...string) *Query {
	for _, val := range orderby {
		v := strings.Split(val, " ")
		v[1] = strings.ToUpper(v[1])
		if v[1] != "ASC" && v[1] != "DESC" {
			panic("invalid sql order type [" + v[1] + "]")
		}
		this._orderby = append(
			this._orderby, SqlIdentifier(v[0]) + " " + v[1])
	}
	return this
}

// limit: integer
func (this *Query) Limit(limit int32) *Query {
	this._limit = limit
	return this
}

// offset: integer
func (this *Query) Offset(offset int32) *Query {
	this._offset = offset
	return this
}

// return all elements matching select query
func (this *Query) All() *[]map[string]any {
	q := "SELECT "
	if len(this._columns) > 0 {
		q += strings.Join(this._columns, ", ")
	} else {
		q += "*"
	}
	q += " FROM " + SqlIdentifier(this._tablename)

	if len(this._filters) > 0 {
		q += "\nWHERE " + strings.Join(this._filters, " ")
	}
	if len(this._groupby) > 0 {
		q += "\nGROUP BY " + strings.Join(this._groupby, ", ")
	}
	if len(this._orderby) > 0 {
		q += "\nORDER BY " + strings.Join(this._orderby, ", ")
	}
	if this._limit >= 0 {
		q += "\nLIMIT " + strconv.Itoa((int)(this._limit))
	}
	if this._offset >= 0 {
		q += "\nOFFSET " + strconv.Itoa((int)(this._offset))
	}
	q += ";"

	this.dbs.Execute(q, this._execargs)
	return this.dbs.FetchAll()
}
