package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Session struct {
	DBh          *Handler
	sqlDB        *sql.DB
	sqlTX        *sql.Tx
	rowsAffected int64
}

func NewSession(dbh *Handler) *Session {
	return &Session{
		DBh: dbh,
	}
}

// create new query for session
func (sess *Session) Query(model IModel) *Query {
	return NewQuery(sess, model)
}

func (sess *Session) IsConnected() bool {
	return sess.sqlDB != nil
}

func (sess *Session) InTransaction() bool {
	return sess.sqlTX != nil
}

func (sess *Session) Connect() error {
	// already connected
	if sess.sqlDB != nil {
		return nil
	}

	if sess.DBh == nil {
		return fmt.Errorf("no database handler defined")
	}
	if sess.DBh.Logger != nil {
		sess.DBh.Logger.Debug(
			"(%v) - connect", sess.DBh.Options["database"])
	}

	sqlDB, err := sess.DBh.Engine.Connect(sess.DBh.Options)
	if err != nil {
		return err
	}
	sess.sqlDB = sqlDB
	return nil
}

func (sess *Session) Close() error {
	// not connected
	if sess.sqlDB == nil {
		return nil
	}
	if sess.DBh.Logger != nil {
		sess.DBh.Logger.Debug(
			"(%v) - closed", sess.DBh.Options["database"])
	}

	sess.sqlDB.Close()
	sess.sqlDB = nil
	sess.sqlTX = nil
	return nil
}

func (sess *Session) Begin() error {
	// connect if not connected
	if err := sess.Connect(); err != nil {
		return err
	}
	if sess.DBh.Logger != nil {
		sess.DBh.Logger.Debug(
			"(%v) - begin", sess.DBh.Options["database"])
	}

	sqlTX, err := sess.sqlDB.Begin()
	if err != nil {
		return err
	}
	sess.sqlTX = sqlTX
	return nil
}

func (sess *Session) Commit() error {
	// not in transaction
	if sess.sqlTX == nil {
		return nil
	}
	if sess.DBh.Logger != nil {
		sess.DBh.Logger.Debug(
			"(%v) - commit", sess.DBh.Options["database"])
	}

	if err := sess.sqlTX.Commit(); err != nil {
		return err
	}
	sess.sqlTX = nil
	return nil
}

func (sess *Session) RollBack() error {
	// not in transaction
	if sess.sqlTX == nil {
		return nil
	}
	if sess.DBh.Logger != nil {
		sess.DBh.Logger.Debug(
			"(%v) - rollback", sess.DBh.Options["database"])
	}

	if err := sess.sqlTX.Rollback(); err != nil {
		return err
	}
	sess.sqlTX = nil
	return nil
}

func (sess *Session) Execute(stmt string, params ...any) error {
	// return error if not connected
	if err := sess.Connect(); err != nil {
		return err
	}

	// format statment args placeholder and log
	stmt = sess.DBh.Engine.FormatSqlStmt(stmt)
	sess.logSql(stmt, params...)

	retries := sess.DBh.Options["retries"].(int)
	retryDelay := int(sess.DBh.Options["retry_delay"].(float64) * 100)

	var err error
	var res sql.Result
	for i := 0; i < retries; i++ {
		if sess.sqlTX != nil {
			res, err = sess.sqlTX.Exec(stmt, params...)
		} else {
			res, err = sess.sqlDB.Exec(stmt, params...)
		}
		// check error
		if err == nil {
			sess.rowsAffected, _ = res.RowsAffected()
			return nil
		} else if sess.isBreakingError(err.Error()) {
			return err
		}
		time.Sleep(time.Millisecond * time.Duration(retryDelay))
	}

	return err
}

func (sess *Session) FetchAll(
	stmt string, params ...any) ([]ModelData, error) {

	// return error if not connected
	if err := sess.Connect(); err != nil {
		return nil, err
	}

	// format statment args placeholder and log
	stmt = sess.DBh.Engine.FormatSqlStmt(stmt)
	sess.logSql(stmt, params...)

	retries := sess.DBh.Options["retries"].(int)
	retryDelay := int(sess.DBh.Options["retry_delay"].(float64) * 100)

	var err error
	var rows *sql.Rows
	var colNames []string

	// query result rows and column names
	for i := 0; i < retries; i++ {
		rows, err = sess.sqlDB.Query(stmt, params...)
		if err == nil {
			defer rows.Close()
			colNames, err = rows.Columns()
			break
		} else if sess.isBreakingError(err.Error()) {
			return nil, err
		}
		time.Sleep(time.Millisecond * time.Duration(retryDelay))
	}
	if err != nil {
		return nil, err
	}

	result := []ModelData{}

	// extract results data
	lenCols := len(colNames)

	// create empty slice to represent cols data, and second
	// slice containing pointers to items in cols data slice.
	colsData := make([]any, lenCols)
	colsPtrs := make([]any, lenCols)
	for k := range colsData {
		colsPtrs[k] = &colsData[k]
	}
	for rows.Next() {
		if err := rows.Scan(colsPtrs...); err != nil {
			return nil, err
		}
		rowData := ModelData{}
		// retrieve value for each column from data slice,
		for k, colName := range colNames {
			rowData[colName] = colsData[k]
		}
		result = append(result, rowData)
	}

	return result, nil
}

func (sess *Session) RowsAffected() int64 {
	return sess.rowsAffected
}

func (sess *Session) logSql(stmt string, params ...any) {
	if sess.DBh.Logger != nil {
		if len(params) > 0 {
			sess.DBh.Logger.Debug(
				"SQL:\n---\n"+stmt+"\nPARAMS: %v\n---", params)
		} else {
			sess.DBh.Logger.Debug(
				"SQL:\n---\n%v\n---", stmt)
		}
	}
}

func (sess *Session) isBreakingError(err string) bool {
	errlist := sess.DBh.Engine.ListRetryErrors()
	for _, match := range errlist {
		if strings.Contains(err, match) {
			return false
		}
	}
	return true
}
