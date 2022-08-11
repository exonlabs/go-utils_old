package sqlite

import (
	"fmt"
	"strings"

	"github.com/exonlabs/go-utils/pkg/db"
)

type Query struct {
	db.BaseQuery
}

func NewQuery(dbs db.ISession, model db.IModel) db.IQuery {
	var this Query
	this.IQuery = &this
	this.DBs = dbs
	this.Model = model
	return &this
}

func (this *Query) CreateTable() {
	var columns, constraints, indexes []string

	for _, c := range this.Model.TableColumns() {
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
					"ON \"%s\" (\"%s\");", u, this.Model.TableName(), c[0], this.Model.TableName(), c[0]))
		}
	}

	defs := columns
	defs = append(defs, constraints...)
	defs = append(defs, this.Model.TableConstraints())

	sql := "CREATE TABLE IF NOT EXISTS " + this.Model.TableName() + " ("
	sql += strings.Join(defs, ",\n")
	sql = strings.TrimSpace(sql)
	sql = strings.TrimSuffix(sql, ",")
	if _, ok := this.Model.TableArgs()["without_rowid"]; ok {
		sql += "\n) WITHOUT ROWID;\n"
	} else {
		sql += "\n);\n"
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
