package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/exonlabs/go-utils/pkg/db"
	"github.com/exonlabs/go-utils/pkg/db/backends/sqlite"
)

var DB_OPTIONS = map[string]any{
	"database": "/tmp/test.db",
}

type foobar struct {
	db.BaseModel
}

func (this *foobar) InitializeData(dbs db.ISession) {
	dbs.Begin()
	for i := 1; i <= 1000; i++ {
		if err := dbs.Query(this).Insert(map[string]any{
			"col1": "foo_" + strconv.Itoa(i),
			"col2": "description_" + strconv.Itoa(i),
			"col3": i,
		}); err != nil {
			dbs.RollBack()
			panic(err)
		}
	}
	dbs.Commit()
}

var Foobar = foobar{
	db.BaseModel{
		Table_Name: "foobar",
		Table_Columns: [][]string{
			{"col1", "TEXT NOT NULL"},
			{"col2", "TEXT"},
			{"col3", "INTEGER"},
			{"col4", "BOOLEAN NOT NULL DEFAULT 0"},
		},
		Table_Constraints: "PRIMARY KEY (\"col1\")",
	},
}

func main() {
	defer panicHandler()
	defer os.Remove(DB_OPTIONS["database"].(string))

	log.Println("DB Options:", DB_OPTIONS)

	dbh := sqlite.NewDBHandler(DB_OPTIONS)
	dbh.InitDatabase([]db.IModel{&Foobar})

	log.Println("DB initialize: Done")

	log.Println("All entries:")
	items := dbh.CreateSession().Query(&Foobar).
		// Columns("col1").
		// FilterBy("col1", 1).
		// Filters("AND (col2=? OR col3 IN (?,?))", true, 2, 3).
		// Filters("AND col3 LIKE ?", "xyz").
		// FilterBy("col99", 99).
		// GroupBy("col1").GroupBy("col2").
		// OrderBy("colx desc", "coly ASC").
		Limit(100).
		// Offset(30).
		Select()

	for _, item := range items {
		fmt.Println(item)
	}

}

func panicHandler() {
	err := recover()
	if err != nil {
		if fmt.Sprint(err) == "EOF" {
			fmt.Printf("\n-- terminated --\n")
			return
		}
		panic(err)
	}
}
