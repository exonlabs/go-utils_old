package main

import (
	"fmt"

	"github.com/exonlabs/go-utils/pkg/db"
	"github.com/exonlabs/go-utils/pkg/db/sqlite"
)

var DB_OPTIONS = map[string]any{
	"database": "/tmp/test.db",
}

// type foobar_model struct {
//     db.BaseModel
// }

// var FoobarModel = foobar_model{

// }

// func test_query() {
//     defer func() {
//         if err := recover(); err != nil {
//             fmt.Println("panic occurred:", err)
//         }
//     }()
//     q := db.NewQuery()
//     q = q.GroupBy("col1", "col2")
// }

func main() {
	defer panicHandler()

	fmt.Printf("\nDB Options:\n%v\n\n", DB_OPTIONS)

	dbs := sqlite.NewSession(nil, nil, 4)

	q := db.NewQuery(dbs, "tbl_name").
		Columns("guid", "col1", "col2", "col3").
		FilterBy("col1", 1).
		Filters("AND (col2=? OR col3 IN (?,?))", true, 2, 3).
		Filters("AND col3 LIKE ?", "xyz").
		FilterBy("col99", 99).
		GroupBy("col1").GroupBy("col2").
		OrderBy("colx desc", "coly ASC").
		Limit(20).Offset(30)

	q.All()
	// fmt.Printf("<%T> = %v\n", res, res)

	// v := "A_s"
	// r, err := db.SqlIdentifier(v)
	// if err != nil {
	//     fmt.Printf("Error: %v\n", err)
	// } else {
	//     fmt.Printf("r = %v\n", r)
	// }

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
