package db

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type Query struct {
	DBs   *Session
	Model IModel

	// runtime table name to use, this allows for mapping
	// same model to multiple tables
	stmtTablename TableName
	// query columns
	stmtColumns []string
	// query sql filters and args
	stmtFilters  []string
	stmtExecargs []any
	// query grouping and ordering
	stmtGroupby []string
	stmtOrderby []string
	stmtHaving  string
	// query data set limit and offset
	stmtLimit  int
	stmtOffset int

	// run time error
	lastError error
}

func NewQuery(dbs *Session, model IModel) *Query {
	qry := Query{
		DBs:           dbs,
		Model:         model,
		stmtTablename: model.GetTableName(),
	}
	if m, ok := model.(IModelDefaultOrders); ok {
		qry.stmtOrderby = m.GetDefaultOrders()
	}
	return &qry
}

// set runtime table name
func (qry *Query) Table(name TableName) *Query {
	qry.stmtTablename = name
	return qry
}

// set columns
func (qry *Query) Columns(columns ...string) *Query {
	qry.stmtColumns = columns
	return qry
}

// add filters
func (qry *Query) Filter(expr string, params ...any) *Query {
	qry.stmtFilters = append(qry.stmtFilters, expr)
	qry.stmtExecargs = append(qry.stmtExecargs, params...)
	return qry
}

// add filters
func (qry *Query) FilterBy(column string, value any) *Query {
	cond := ""
	if len(qry.stmtFilters) > 0 {
		cond = "AND "
	}
	qry.stmtFilters = append(
		qry.stmtFilters, cond+column+"="+SQL_PLACEHOLDER)
	qry.stmtExecargs = append(qry.stmtExecargs, value)
	return qry
}

// add grouping
func (qry *Query) GroupBy(groupby ...string) *Query {
	qry.stmtGroupby = groupby
	return qry
}

// add ordering: "colname ASC|DESC"
func (qry *Query) OrderBy(orderby ...string) *Query {
	qry.stmtOrderby = orderby
	return qry
}

// add having expr
func (qry *Query) Having(expr string, val any) *Query {
	qry.stmtHaving = expr
	qry.stmtExecargs = append(qry.stmtExecargs, val)
	return qry
}

// add limit
func (qry *Query) Limit(limit int) *Query {
	qry.stmtLimit = limit
	return qry
}

// add offset
func (qry *Query) Offset(offset int) *Query {
	qry.stmtOffset = offset
	return qry
}

// return all elements matching select query
func (qry *Query) All() ([]ModelData, error) {
	if qry.lastError != nil {
		return nil, qry.lastError
	}

	limitPrefix := ""
	if qry.DBs.DBh.Engine.GetBackendName() == "mssql" {
		if qry.stmtLimit > 0 && len(qry.stmtOrderby) == 0 {
			limitPrefix = fmt.Sprintf("TOP(%v) ", qry.stmtLimit)
		}
	}

	q := "SELECT " + limitPrefix
	if len(qry.stmtColumns) > 0 {
		q += strings.Join(qry.stmtColumns, ", ")
	} else {
		q += "*"
	}
	q += " FROM " + qry.stmtTablename

	if len(qry.stmtFilters) > 0 {
		q += "\nWHERE " + strings.Join(qry.stmtFilters, " ")
	}
	if len(qry.stmtGroupby) > 0 {
		q += "\nGROUP BY " + strings.Join(qry.stmtGroupby, ", ")
	}
	if len(qry.stmtHaving) > 0 {
		q += "\nHAVING " + qry.stmtHaving
	}
	if len(qry.stmtOrderby) > 0 {
		q += "\nORDER BY " + strings.Join(qry.stmtOrderby, ", ")
	}
	if qry.DBs.DBh.Engine.GetBackendName() == "mssql" {
		if qry.stmtOffset > 0 || qry.stmtLimit > 0 {
			q += fmt.Sprintf(
				"\nOFFSET %v ROWS", qry.stmtOffset)
		}
		if qry.stmtLimit > 0 {
			q += fmt.Sprintf(
				"\nFETCH NEXT %v ROWS ONLY", qry.stmtLimit)
		}
	} else {
		if qry.stmtLimit > 0 {
			q += fmt.Sprintf("\nLIMIT %v", qry.stmtLimit)
		}
		if qry.stmtOffset > 0 {
			q += fmt.Sprintf("\nOFFSET %v", qry.stmtOffset)
		}
	}
	q += ";"

	// run query and fetch data
	result, err := qry.DBs.FetchAll(q, qry.stmtExecargs...)
	if err != nil {
		return nil, err
	}

	// if we have data converters
	if m, ok := qry.Model.(IModelDataConverters); ok {
		var res []ModelData
		for _, d := range result {
			res = append(res, qry.dataMapper(d, m.GetDataConverters()))
		}
		return res, nil
	}

	return result, nil
}

// return first element matching select query
func (qry *Query) First() (ModelData, error) {
	qry.stmtLimit, qry.stmtOffset = 1, 0
	result, err := qry.All()
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return nil, nil
}

// return one element matching filter params or nil
// there must be only one element or none
func (qry *Query) One() (ModelData, error) {
	qry.stmtLimit, qry.stmtOffset = 2, 0
	result, err := qry.All()
	if err != nil {
		return nil, err
	}
	switch len(result) {
	case 0:
		return nil, nil
	case 1:
		return result[0], nil
	}
	return nil, fmt.Errorf("multiple entries found")
}

// get model data defined by guid
func (qry *Query) Get(guid string) (ModelData, error) {
	qry.stmtFilters = []string{"guid=" + SQL_PLACEHOLDER}
	qry.stmtExecargs = []any{guid}
	qry.stmtGroupby = []string{}
	qry.stmtOrderby = []string{}
	qry.stmtHaving = ""
	return qry.One()
}

func (qry *Query) Count() (int64, error) {
	if qry.lastError != nil {
		return 0, qry.lastError
	}

	q := "SELECT count(*) as count FROM " + qry.stmtTablename
	if len(qry.stmtFilters) > 0 {
		q += "\nWHERE " + strings.Join(qry.stmtFilters, " ")
	}
	if len(qry.stmtGroupby) > 0 {
		q += "\nGROUP BY " + strings.Join(qry.stmtGroupby, ", ")
	}
	q += ";"

	data, err := qry.DBs.FetchAll(q, qry.stmtExecargs...)
	if err != nil {
		return 0, err
	}

	return data[0]["count"].(int64), nil
}

func (qry *Query) Insert(data ModelData) (string, error) {
	if qry.lastError != nil {
		return "", qry.lastError
	}

	// check and create guid in data
	if data["guid"] == nil || len(data["guid"].(string)) == 0 {
		data["guid"] = qry.GenerateGuid()
	}

	// apply data adapters
	if m, ok := qry.Model.(IModelDataAdapters); ok {
		data = qry.dataMapper(data, m.GetDataAdapters())
	}

	columns, holders, execArgs := []string{}, []string{}, []any{}
	for k, v := range data {
		columns = append(columns, k)
		holders = append(holders, SQL_PLACEHOLDER)
		execArgs = append(execArgs, v)
	}

	q := "INSERT INTO " + qry.stmtTablename
	q += fmt.Sprintf("\n(%v)", strings.Join(columns, ", "))
	q += fmt.Sprintf("\nVALUES (%v)", strings.Join(holders, ", "))
	q += ";"

	err := qry.DBs.Execute(q, execArgs...)
	if err != nil {
		return "", err
	}

	return data["guid"].(string), nil
}

func (qry *Query) Update(data ModelData) (int64, error) {
	if qry.lastError != nil {
		return 0, qry.lastError
	}

	// apply data adapters
	if m, ok := qry.Model.(IModelDataAdapters); ok {
		data = qry.dataMapper(data, m.GetDataAdapters())
	}

	columns, execArgs := []string{}, []any{}
	for k, v := range data {
		columns = append(columns, k+"="+SQL_PLACEHOLDER)
		execArgs = append(execArgs, v)
	}

	q := "UPDATE " + qry.stmtTablename
	q += "\nSET " + strings.Join(columns, ", ")
	if len(qry.stmtFilters) > 0 {
		q += "\nWHERE " + strings.Join(qry.stmtFilters, " ")
	}

	execArgs = append(execArgs, qry.stmtExecargs...)

	err := qry.DBs.Execute(q, execArgs...)
	if err != nil {
		return 0, err
	}

	return qry.DBs.RowsAffected(), nil
}

func (qry *Query) Delete() (int64, error) {
	if qry.lastError != nil {
		return 0, qry.lastError
	}

	q := "DELETE FROM " + qry.stmtTablename
	if len(qry.stmtFilters) > 0 {
		q += "\nWHERE " + strings.Join(qry.stmtFilters, " ")
	}

	err := qry.DBs.Execute(q, qry.stmtExecargs...)
	if err != nil {
		return 0, err
	}

	return qry.DBs.RowsAffected(), nil
}

// check valid sql identifier string
func (qry *Query) SqlIdent(name string) string {
	if qry.lastError != nil {
		return ""
	}

	match, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", name)
	if match {
		return name
	}

	qry.lastError = fmt.Errorf(
		fmt.Sprintf("invalid sql identifier '%v'", name))
	return ""
}

func (qry *Query) GenerateGuid() string {
	return hex.EncodeToString(
		uuid.NewV5(uuid.NewV1(), uuid.NewV4().String()).Bytes())
}

func (qry *Query) dataMapper(
	data ModelData, mappers ModelDataMappers) ModelData {
	for key, fn := range mappers {
		if val, ok := data[key]; ok {
			data[key] = fn(val)
		}
	}
	return data
}
