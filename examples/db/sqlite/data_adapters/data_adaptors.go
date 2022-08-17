package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

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
			"col5": []string{"foo_" + strconv.Itoa(i),
				"description_" + strconv.Itoa(i),
				"pass_" + strconv.Itoa(i),
			},
			"password": "pass_" + strconv.Itoa(i),
		}); err != nil {
			dbs.RollBack()
			panic(err)
		}
	}
	dbs.Commit()
}

var Foobar = &foobar{
	db.BaseModel{
		TableName: "foobar",
		TableColumns: [][]string{
			{"col1", "TEXT NOT NULL", "UNIQUE INDEX"},
			{"col2", "TEXT"},
			{"col3", "INTEGER"},
			{"col4", "BOOLEAN NOT NULL DEFAULT 0"},
			{"col5", "TEXT NOT NULL"},
			{"password", "TEXT NOT NULL"},
		},
		DataAdapters: map[string]func(any) any{
			"col5": func(slice any) any {
				return strings.Join(slice.([]string), ", ")
			},
			"password": func(text any) any {
				data := fmt.Sprint(text)
				hash := sha256.Sum256([]byte(data))
				return hex.EncodeToString(hash[:])
			},
		},
		DataConverters: map[string]func(any) any{
			"col5": func(text any) any {
				return strings.Split(text.(string), ", ")
			},
		},
	},
}

func main() {
	defer panicHandler()
	defer os.Remove(DB_OPTIONS["database"].(string))

	fmt.Println("DB Options:", DB_OPTIONS)

	dbh := sqlite.NewDBHandler(DB_OPTIONS)
	dbh.InitializeDB([]db.IModel{Foobar})

	fmt.Println("DB initialize: Done")

	var items []map[string]any

	fmt.Println("\nGet 7 entries:")
	items = dbh.Session().Query(Foobar).
		Limit(7).Select()
	for _, item := range items {
		fmt.Println(item)
	}
	fmt.Println("Total Items:", dbh.Session().Query(Foobar).Count())

	fmt.Println("")
}

func panicHandler() {
	err := recover()
	if err != nil {
		panic(err)
	}
}
