package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/exonlabs/go-utils_old/pkg/db"
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
	return "mysql"
}

func (eng *Engine) FormatSqlStmt(stmt string) string {
	return strings.Replace(
		stmt, db.SQL_PLACEHOLDER, "?", -1)
}

func (eng *Engine) Connect(options KwArgs) (*sql.DB, error) {
	// params
	database, _ := options["database"].(string)
	if len(database) == 0 {
		return nil, fmt.Errorf("invalid database configuration")
	}
	host, _ := options["host"].(string)
	if len(host) == 0 {
		host = "localhost"
	}
	port, _ := options["port"].(int)
	if port <= 0 {
		port = 3306
	}
	username, _ := options["username"].(string)
	password, _ := options["password"].(string)
	extargs, _ := options["extargs"].(string)
	if !strings.Contains(extargs, "collation=") {
		extargs = "collation=utf8mb4_unicode_ci&" + extargs
	}
	if !strings.Contains(extargs, "charset=") {
		extargs = "charset=utf8mb4,utf8&" + extargs
	}

	// create data source name
	dsn := fmt.Sprintf("tcp(%v:%d)/%v?%v", host, port, database, extargs)
	if len(username) > 0 {
		if len(password) > 0 {
			dsn = fmt.Sprintf("%v:%v@%v", username, password, dsn)
		} else {
			dsn = fmt.Sprintf("%v@%v", username, dsn)
		}
	}

	sqlDB, err := sql.Open("mysql", dsn)
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
	// 		if strings.Contains(c[1], "0") {
	// 			c[1] = strings.Replace(c[1], "0", "false", -1)
	// 		} else {
	// 			c[1] = strings.Replace(c[1], "1", "true", -1)
	// 		}
	// 		expr = append(expr, db.SqlIdentifier(c[0], &err)+" "+c[1])
	// 		if err != nil {
	// 			return nil, fmt.Errorf("table columns: " + err.Error())
	// 		}
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
	// 			"CREATE %vINDEX IF NOT EXISTS "+
	// 				"ix_%v_%v "+
	// 				"ON \"%v\" (\"%v\");", u,
	// 			tblname, c[0], tblname, c[0]))
	// 	}
	// }

	// expr = append(expr, constraints...)
	// expr = append(expr, model.TableConstraints()...)

	// stmt := "CREATE TABLE IF NOT EXISTS \"" + tblname + "\" (\n"
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
