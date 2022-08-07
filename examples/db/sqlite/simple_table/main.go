package main

import (
    "fmt"

    "github.com/exonlabs/go-utils/pkg/db"
    "github.com/exonlabs/go-utils/pkg/db/backends/sqlite"
)

var DB_OPTIONS = map[string]any{
    "database": "/tmp/test.db",
}

var Foobar := db.BaseModel{
    TableName: "foobar",
    TableColumns: [][]string{
        {"col1", "TEXT NOT NULL"},
        {"col2", "TEXT"},
        {"col3", "INTEGER"},
        {"col4", "BOOLEAN NOT NULL DEFAULT 0"},
    },
    TableConstraints: "PRIMARY KEY (\"col1\")",
    InitializeData: func(model db.BaseModel, dbs db.Session) {
        for i := 0; i < 5; i++ {
            dbs.Query(model).Model.Create(dbs, map[string]any{
                "col1": "foo_" + strconv.Itoa(i),
                "col2": "description_" + strconv.Itoa(i),
                "col3": i,
            })
        }
    },
}


func main() {
    defer panicHandler()

    fmt.Printf("\nDB Options:\n%v\n\n", DB_OPTIONS)






}

func panicHandler() {
    err := recover();
    if err != nil {
        if fmt.Sprint(err) == "EOF" {
            fmt.Printf("\n-- terminated --\n")
            return
        }
        panic(err)
    }
}
