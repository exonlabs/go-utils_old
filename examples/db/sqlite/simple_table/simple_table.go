package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/exonlabs/go-utils/pkg/db"
	"github.com/exonlabs/go-utils/pkg/db/backends/sqlite"
	"github.com/exonlabs/go-utils/pkg/logging"
)

var DB_OPTIONS = map[string]any{
	"database": "/tmp/test.db",
	// "sql_argmap": "$?",
}

type foobar struct {
	db.BaseModel
}

func (this *foobar) InitializeData(dbs db.ISession) {
	for i := 1; i <= 5; i++ {
		if err := dbs.Query(this).Insert(map[string]any{
			"col1": "foo_" + strconv.Itoa(i),
			"col2": "description_" + strconv.Itoa(i),
			"col3": i,
		}); err != nil {
			panic(err)
		}
	}
}

var Foobar = &foobar{
	db.BaseModel{
		TableName: "foobar",
		TableColumns: [][]string{
			{"col1", "TEXT NOT NULL", "UNIQUE INDEX"},
			{"col2", "TEXT"},
			{"col3", "INTEGER"},
			{"col4", "BOOLEAN NOT NULL DEFAULT 0"},
		},
	},
}

func main() {
	defer panicHandler()
	defer os.Remove(DB_OPTIONS["database"].(string))

	fmt.Println("DB Options:", DB_OPTIONS)

	dbh := sqlite.NewDBHandler(DB_OPTIONS)
	dbh.Logger = logging.NewLogger()
	if len(os.Args[1:]) > 0 {
		if os.Args[1:][0] == "-x" {
			dbh.Logger.Level = logging.DEBUG
		}
	}
	dbh.InitializeDB([]db.IModel{Foobar})

	fmt.Println("DB initialize: Done")

	var items []map[string]any

	fmt.Println("\nGet all entries:")
	items = dbh.Session().Query(Foobar).Select()
	for _, item := range items {
		fmt.Println(item)
	}
	fmt.Println("Total Items:", dbh.Session().Query(Foobar).Count())

	fmt.Println("\nGet custom columns entries:")
	items = dbh.Session().Query(Foobar).
		Columns("col1").Limit(2).OrderBy("col1 ASC").Select()
	for _, item := range items {
		fmt.Println(item)
	}

	fmt.Println("\nGet filter columns entries:")
	items = dbh.Session().Query(Foobar).
		Filters("(col2=$? OR col3 IN ($?,$?))", "description_3", 1, 3).
		OrderBy("col1 ASC").Select()
	for _, item := range items {
		fmt.Println(item)
	}

	fmt.Println("\nUpdate first row")
	dbh.Session().Query(Foobar).
		FilterBy("col3", 1).
		Update(map[string]any{"col1": "boo_1", "col2": "boo_2"})
	fmt.Println("\nGet all entries:")
	items = dbh.Session().Query(Foobar).Select()
	for _, item := range items {
		fmt.Println(item)
	}

	fmt.Println("\nDELETE second row")
	dbh.Session().Query(Foobar).
		FilterBy("col3", 2).Delete()
	fmt.Println("\nGet all entries:")
	items = dbh.Session().Query(Foobar).Select()
	for _, item := range items {
		fmt.Println(item)
	}

	fmt.Println("")
}

func panicHandler() {
	err := recover()
	if err != nil {
		panic(err)
	}
}
