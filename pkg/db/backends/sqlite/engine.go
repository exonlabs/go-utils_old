package sqlite

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/exonlabs/go-utils_old/pkg/db"
	"github.com/mattn/go-sqlite3"
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
	return "sqlite"
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
	extargs, _ := options["extargs"].(string)
	if !strings.Contains(extargs, "_foreign_keys=") {
		extargs = "_foreign_keys=1&" + extargs
	}

	// create data source name
	dsn := fmt.Sprintf("%v?%v", database, extargs)

	sqlDB, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	return sqlDB, nil
}

func (eng *Engine) GenTableSchema(
	tblname db.TableName, meta db.TableMeta) ([]string, error) {

	tblcolumns := meta.Columns
	// add guid column if not exist as first column
	if tblcolumns[0][0] != "guid" {
		tblcolumns = append(db.TableColumns{
			{"guid", "VARCHAR(32) NOT NULL", "PRIMARY"},
		}, tblcolumns...)
	}

	var expr, constraints, indexes []string
	for _, c := range tblcolumns {
		expr = append(expr, c[0]+" "+c[1])

		// add check constraint for bool type
		if strings.Contains(c[1], "BOOLEAN") {
			constraints = append(constraints,
				fmt.Sprintf("CHECK (\"%v\" IN (0,1))", c[0]))
		}

		if len(c) <= 2 {
			continue
		}

		if strings.Contains(c[2], "PRIMARY") {
			// add primary_key constraint
			constraints = append(constraints,
				fmt.Sprintf("PRIMARY KEY (\"%v\")", c[0]))
		} else if strings.Contains(c[2], "UNIQUE") &&
			!strings.Contains(c[2], "INDEX") {
			// add unique constraint if not indexed column
			constraints = append(constraints,
				fmt.Sprintf("UNIQUE (\"%v\")", c[0]))
		}

		if strings.Contains(c[2], "PRIMARY") ||
			strings.Contains(c[2], "INDEX") {

			u := ""
			if strings.Contains(c[2], "PRIMARY") ||
				strings.Contains(c[2], "UNIQUE") {
				u = "UNIQUE "
			}

			indexes = append(indexes, fmt.Sprintf(
				"CREATE %vINDEX IF NOT EXISTS "+
					"ix_%v_%v "+
					"ON \"%v\" (\"%v\");", u,
				tblname, c[0], tblname, c[0]))
		}
	}
	expr = append(expr, constraints...)

	// add explicit table constraints
	for _, c := range meta.Constraints {
		expr = append(expr,
			fmt.Sprintf("CONSTRAINT %v %v", c[0], c[1]))
	}

	sql := "CREATE TABLE IF NOT EXISTS \"" + tblname + "\" (\n"
	sql += strings.Join(expr, ",\n")
	sql = strings.TrimSpace(sql)
	sql = strings.TrimSuffix(sql, ",")

	val, ok := meta.Options["without_rowid"]
	if ok && !val.(bool) {
		sql += "\n);"
	} else {
		sql += "\n) WITHOUT ROWID;"
	}

	result := []string{sql}
	result = append(result, indexes...)
	return result, nil
}

func (eng *Engine) ListRetryErrors() []string {
	return []string{
		// The database file is locked
		sqlite3.ErrBusy.Error(),
		// A table in the database is locked
		sqlite3.ErrLocked.Error(),
		// Some kind of disk I/O error occurred
		sqlite3.ErrIoErr.Error(),
		// Unable to open the database file
		sqlite3.ErrCantOpen.Error(),
	}
}
