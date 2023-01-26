package mssql

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/exonlabs/go-utils/db"
)

type KwArgs = map[string]any

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func NewHandler(options KwArgs) *db.Handler {
	return db.NewHandler(NewEngine(), options)
}

func (eng *Engine) GetBackendName() string {
	return "mssql"
}

func (eng *Engine) FormatSqlStmt(stmt string) string {
	return strings.Replace(
		stmt, db.SQL_PLACEHOLDER, "?", -1)
}

func (eng *Engine) Connect(options KwArgs) (*sql.DB, error) {
	for _, key := range []string{
		"database", "host", "username", "password"} {
		if v, ok := options[key].(string); !ok || len(v) == 0 {
			return nil, fmt.Errorf("invalid database configuration")
		}
	}
	if v, ok := options["port"].(int); !ok || v == 0 {
		return nil, fmt.Errorf("invalid database configuration")
	}

	uri := fmt.Sprintf(
		"server=%v;user id=%v;password=%v;port=%d;database=%v",
		options["host"], options["username"], options["password"],
		options["port"], options["database"])

	sqlDB, err := sql.Open("mssql", uri)
	if err != nil {
		return nil, err
	}
	return sqlDB, nil
}

func (eng *Engine) GenTableSchema(
	tblname db.TableName, meta db.TableMeta) ([]string, error) {

	// var err error

	// tblname := model.TableName()
	// if kwargs != nil {
	// 	if val, ok := kwargs["table_name"]; ok {
	// 		tblname = val.(string)
	// 	}
	// }
	// tblname = db.SqlIdentifier(tblname, &err)
	// if err != nil {
	// 	return nil, fmt.Errorf("table name: " + err.Error())
	// }

	// tblcolumns := model.TableColumns()
	// if tblcolumns[0][0] != "guid" {
	// 	tblcolumns = append(db.TColumns{
	// 		{"guid", "VARCHAR(32) NOT NULL", "PRIMARY"},
	// 	}, tblcolumns...)
	// }

	// var expr, constraints, indexes []string
	// for _, c := range tblcolumns {
	// 	if strings.Contains(c[1], "BOOLEAN") {
	// 		c[1] = strings.Replace(c[1], "BOOLEAN", "BIT", -1)
	// 		expr = append(expr, db.SqlIdentifier(c[0], &err)+" "+c[1])
	// 		if err != nil {
	// 			return nil, fmt.Errorf("table columns: " + err.Error())
	// 		}
	// 		constraints = append(constraints,
	// 			fmt.Sprintf("CHECK (\"%v\" IN (0,1))", c[0]))
	// 	} else {
	// 		expr = append(expr, db.SqlIdentifier(c[0], &err)+" "+c[1])
	// 		if err != nil {
	// 			return nil, fmt.Errorf("table columns: " + err.Error())
	// 		}
	// 	}

	// 	if len(c) <= 2 {
	// 		continue
	// 	}

	// 	if strings.Contains(c[2], "PRIMARY") {
	// 		constraints = append(constraints,
	// 			fmt.Sprintf("PRIMARY KEY (\"%v\")", c[0]))
	// 	} else if strings.Contains(c[2], "UNIQUE") && !strings.Contains(c[2], "INDEX") {
	// 		constraints = append(constraints,
	// 			fmt.Sprintf("UNIQUE (\"%v\")", c[0]))
	// 	}

	// 	if strings.Contains(c[2], "PRIMARY") || strings.Contains(c[2], "INDEX") {
	// 		u := "UNIQUE "
	// 		indexes = append(indexes, fmt.Sprintf(
	// 			"IF NOT EXISTS (SELECT * FROM sys.indexes "+
	// 				"WHERE name='ix_%v_%v')\n"+
	// 				"CREATE %vINDEX "+
	// 				"ix_%v_%v "+
	// 				"ON %v (%v);",
	// 			tblname, c[0], u, tblname, c[0], tblname, c[0]))
	// 	}
	// }

	// expr = append(expr, constraints...)
	// expr = append(expr, model.TableConstraints()...)

	// stmt := "IF OBJECT_ID(N'" + tblname + "','U') IS NULL\n"
	// stmt += "CREATE TABLE \"" + tblname + "\" (\n"
	// stmt += strings.Join(expr, ",\n")
	// stmt = strings.TrimSpace(stmt)
	// stmt = strings.TrimSuffix(stmt, ",")
	// stmt += "\n);"

	// result := []string{stmt}
	// result = append(result, indexes...)
	// return result, nil
	return []string{}, nil
}

func (eng *Engine) ListRetryErrors() []string {
	return []string{}
}
