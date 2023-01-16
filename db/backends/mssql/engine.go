package mssql

// import (
// 	"database/sql"
// 	"fmt"
// 	"strings"

// 	_ "github.com/denisenkom/go-mssqldb"

// 	"github.com/exonlabs/go-utils/pkg/db"
// 	. "github.com/exonlabs/go-utils/pkg/globals"
// )

// type Engine struct{}

// func NewEngine() *Engine {
// 	return &Engine{}
// }

// func (this *Engine) BackendName() string {
// 	return "mssql"
// }

// func (this *Engine) SqlDB(options TArgs) (*sql.DB, error) {
// 	for _, v := range []string{"host", "user", "password", "port", "database"} {
// 		_, ok := options[v]
// 		if !ok {
// 			panic("invalid database configuration")
// 		}
// 	}

// 	connString := fmt.Sprintf("server=%v;user id=%v;password=%v;port=%d;database=%v",
// 		options["host"], options["user"], options["password"],
// 		options["port"], options["database"])

// 	db, err := sql.Open("mssql", connString)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return db, nil
// }

// func (this *Engine) PostConnect(db *sql.DB, options TArgs) error {
// 	return nil
// }

// func (this *Engine) TableSchema(model db.IModel, kwargs TArgs) ([]string, error) {
// 	var err error

// 	tblname := model.TableName()
// 	if kwargs != nil {
// 		if val, ok := kwargs["table_name"]; ok {
// 			tblname = val.(string)
// 		}
// 	}
// 	tblname = db.SqlIdentifier(tblname, &err)
// 	if err != nil {
// 		return nil, fmt.Errorf("table name: " + err.Error())
// 	}

// 	tblcolumns := model.TableColumns()
// 	if tblcolumns[0][0] != "guid" {
// 		tblcolumns = append(db.TColumns{
// 			{"guid", "VARCHAR(32) NOT NULL", "PRIMARY"},
// 		}, tblcolumns...)
// 	}

// 	var expr, constraints, indexes []string
// 	for _, c := range tblcolumns {
// 		if strings.Contains(c[1], "BOOLEAN") {
// 			c[1] = strings.Replace(c[1], "BOOLEAN", "BIT", -1)
// 			expr = append(expr, db.SqlIdentifier(c[0], &err)+" "+c[1])
// 			if err != nil {
// 				return nil, fmt.Errorf("table columns: " + err.Error())
// 			}
// 			constraints = append(constraints,
// 				fmt.Sprintf("CHECK (\"%v\" IN (0,1))", c[0]))
// 		} else {
// 			expr = append(expr, db.SqlIdentifier(c[0], &err)+" "+c[1])
// 			if err != nil {
// 				return nil, fmt.Errorf("table columns: " + err.Error())
// 			}
// 		}

// 		if len(c) <= 2 {
// 			continue
// 		}

// 		if strings.Contains(c[2], "PRIMARY") {
// 			constraints = append(constraints,
// 				fmt.Sprintf("PRIMARY KEY (\"%v\")", c[0]))
// 		} else if strings.Contains(c[2], "UNIQUE") && !strings.Contains(c[2], "INDEX") {
// 			constraints = append(constraints,
// 				fmt.Sprintf("UNIQUE (\"%v\")", c[0]))
// 		}

// 		if strings.Contains(c[2], "PRIMARY") || strings.Contains(c[2], "INDEX") {
// 			u := "UNIQUE "
// 			indexes = append(indexes, fmt.Sprintf(
// 				"IF NOT EXISTS (SELECT * FROM sys.indexes "+
// 					"WHERE name='ix_%v_%v')\n"+
// 					"CREATE %vINDEX "+
// 					"ix_%v_%v "+
// 					"ON %v (%v);",
// 				tblname, c[0], u, tblname, c[0], tblname, c[0]))
// 		}
// 	}

// 	expr = append(expr, constraints...)
// 	expr = append(expr, model.TableConstraints()...)

// 	stmt := "IF OBJECT_ID(N'" + tblname + "','U') IS NULL\n"
// 	stmt += "CREATE TABLE \"" + tblname + "\" (\n"
// 	stmt += strings.Join(expr, ",\n")
// 	stmt = strings.TrimSpace(stmt)
// 	stmt = strings.TrimSuffix(stmt, ",")
// 	stmt += "\n);"

// 	result := []string{stmt}
// 	result = append(result, indexes...)

// 	return result, nil
// }

// func (this *Engine) SqlStmtMapper(sql string, options TArgs) string {
// 	return strings.Replace(
// 		sql, options["sql_placeholder"].(string), "?", -1)
// }

// func (this *Engine) DatabaseErrors() []string {
// 	return []string{}
// }
