package sqlite

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/exonlabs/go-utils/pkg/db"
)

type Query struct {
	db.BaseQuery
}

func NewQuery(dbs *Session, model db.IModel) *Query {
	var this Query
	this.IQuery = &this
	this.DBs = dbs
	this.Model = model
	this.StTablename = reflect.Indirect(reflect.ValueOf(model)).
		FieldByName("TableName").String()

	if val, ok := dbs.Options["sql_argmap"]; ok {
		this.SqlArgsMap = val.(string)
	} else {
		this.SqlArgsMap = "$?"
	}

	return &this
}

func (this *Query) CreateTable() {
	var columns, constraints, indexes []string

	model := reflect.Indirect(reflect.ValueOf(this.Model)).
		FieldByName("BaseModel").Interface().(db.BaseModel)

	if model.TableColumns[0][0] != "guid" {
		columns = append(columns, "guid TEXT NOT NULL")
		constraints = append(constraints, "PRIMARY KEY (guid)")
		indexes = append(indexes, "CREATE UNIQUE INDEX ix_"+model.TableName+"_guid ON "+model.TableName+" (guid)")
	}

	for _, c := range model.TableColumns {
		columns = append(columns, db.SqlIdentifier(c[0])+" "+c[1])

		if strings.Contains(c[1], "BOOLEAN") {
			constraints = append(constraints, fmt.Sprintf("CHECK (\"%s\" IN (0,1))", c[0]))
		}

		if len(c) <= 2 {
			continue
		}

		if strings.Contains(c[2], "PRIMARY") {
			constraints = append(constraints, fmt.Sprintf("PRIMARY KEY (\"%s\")", c[0]))
		} else if strings.Contains(c[2], "UNIQUE") && !strings.Contains(c[2], "INDEX") {
			constraints = append(constraints, fmt.Sprintf("UNIQUE (\"%s\")", c[0]))
		}

		if strings.Contains(c[2], "PRIMARY") || strings.Contains(c[2], "INDEX") {
			u := "UNIQUE "
			indexes = append(indexes, fmt.Sprintf(
				"CREATE %sINDEX IF NOT EXISTS "+
					"ix_%s_%s "+
					"ON \"%s\" (\"%s\");", u, this.StTablename, c[0], this.StTablename, c[0]))
		}
	}

	defs := columns
	defs = append(defs, constraints...)
	defs = append(defs, model.TableConstraints)

	sql := "CREATE TABLE IF NOT EXISTS " + this.StTablename + " ("
	sql += strings.Join(defs, ",\n")
	sql = strings.TrimSpace(sql)
	sql = strings.TrimSuffix(sql, ",")
	if _, ok := model.TableArgs["without_rowid"]; ok {
		sql += "\n) WITHOUT ROWID;\n"
	} else {
		sql += "\n);"
	}

	// open transaction
	this.DBs.Begin()
	if err := this.DBs.Execute(sql); err != nil {
		this.DBs.RollBack()
		panic(err)
	}
	// create indexes
	for _, sql := range indexes {
		if err := this.DBs.Execute(sql); err != nil {
			this.DBs.RollBack()
			panic(err)
		}
	}

	this.DBs.Commit()
}
