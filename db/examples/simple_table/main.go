package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/exonlabs/go-utils/db"
	"github.com/exonlabs/go-utils/db/backends/mssql"
	"github.com/exonlabs/go-utils/db/backends/mysql"
	"github.com/exonlabs/go-utils/db/backends/pgsql"
	"github.com/exonlabs/go-utils/db/backends/sqlite"
	"github.com/exonlabs/go-utils/logging"
	"github.com/exonlabs/go-utils/logging/handlers"
	"golang.org/x/exp/maps"
)

type KwArgs = map[string]any

var (
	BACKENDS = []string{"sqlite", "mysql", "pgsql", "mssql"}

	DB_OPTIONS = KwArgs{
		"database": "test",
		"host":     "localhost",
		"port":     0,
		"username": "user1",
		"password": "123456",
		"extargs":  "",

		// optional
		// "connect_timeout": 30,
		// "retries": 10,
		// "retry_delay": 0.5,
	}
)

type Foobar db.BaseModel

func (mdl *Foobar) GetTableName() db.TableName {
	return "foobar"
}

func (mdl *Foobar) GetTableMeta() db.TableMeta {
	return db.TableMeta{
		Columns: db.TableColumns{
			{"col1", "VARCHAR(128) NOT NULL", "UNIQUE INDEX"},
			{"col2", "TEXT"},
			{"col3", "INTEGER"},
			{"col4", "BOOLEAN NOT NULL DEFAULT 0"},
		},
	}
}

func (mdl *Foobar) GetDefaultOrders() []string {
	return []string{"col1 ASC"}
}

func (mdl *Foobar) InitializeData(
	dbs *db.Session, tblname db.TableName) error {
	for i := 0; i < 5; i++ {
		num, _ := dbs.Query(mdl).Table(tblname).
			Filter("col1=$?", "foo_"+strconv.Itoa(i)).Count()
		if num > 0 {
			continue
		}
		_, err := dbs.Query(mdl).Table(tblname).Insert(db.ModelData{
			"col1": "foo_" + strconv.Itoa(i),
			"col2": "description_" + strconv.Itoa(i),
			"col3": i,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func PrintModelData(data []db.ModelData) {
	if len(data) > 0 {
		keys := maps.Keys(data[0])
		sort.Strings(keys)
		for _, item := range data {
			res := []string{}
			for _, k := range keys {
				res = append(res, fmt.Sprintf("%v: %v", k, item[k]))
			}
			fmt.Println(strings.Join(res, ", "))
		}
	}
}

func runOperations(dbh *db.Handler) {
	// define tables
	tables := map[db.TableName]db.IModel{
		"foobar": &Foobar{},
	}

	if err := dbh.InitDatabase(tables); err != nil {
		fmt.Println("ERROR:", err.Error())
		return
	}
	fmt.Println("\nDB initialize: Done")

	dbs := dbh.Session()
	defer dbs.Close()

	// checking DB
	fmt.Println("\nAll entries:")
	if items, err := dbs.Query(&Foobar{}).All(); err != nil {
		fmt.Println("ERROR:", err.Error())
		return
	} else {
		PrintModelData(items)
	}
	if total, err := dbs.Query(&Foobar{}).Count(); err != nil {
		fmt.Println("ERROR:", err.Error())
		return
	} else {
		fmt.Println("\nTotal Items:", total)
	}

	// custom columns
	fmt.Println("\nGet custom columns entries:")
	if items, err := dbs.Query(&Foobar{}).Columns("col1", "col2").
		Limit(2).OrderBy("col1 DESC").All(); err != nil {
		fmt.Println("ERROR:", err.Error())
		return
	} else {
		PrintModelData(items)
	}

	// filtered entries
	fmt.Println("\nGet filter columns entries:")
	if items, err := dbs.Query(&Foobar{}).
		Filter("col2 LIKE $? OR col3 IN ($?,$?)", "description_3", 1, 3).
		All(); err != nil {
		fmt.Println("ERROR:", err.Error())
		return
	} else {
		PrintModelData(items)
	}

	// update entries
	fmt.Println("\nModify: first row")
	if _, err := dbs.Query(&Foobar{}).FilterBy("col3", 1).
		Update(db.ModelData{
			"col1": "boo_1", "col2": "boo_2", "col4": 1,
		}); err != nil {
		fmt.Println("ERROR:", err.Error())
		return
	}
	fmt.Println("-- After Modify --")
	if items, err := dbs.Query(&Foobar{}).All(); err != nil {
		fmt.Println("ERROR:", err.Error())
		return
	} else {
		PrintModelData(items)
	}

	fmt.Println("\nDelete: first row")
	if _, err := dbs.Query(&Foobar{}).FilterBy("col3", 1).
		Delete(); err != nil {
		fmt.Println("ERROR:", err.Error())
		return
	}
	fmt.Println("-- After Delete --")
	if items, err := dbs.Query(&Foobar{}).All(); err != nil {
		fmt.Println("ERROR:", err.Error())
		return
	} else {
		PrintModelData(items)
	}
}

func main() {
	logger := handlers.NewStdoutLogger("main")
	logger.Level = logging.LEVEL_INFO
	logger.Formatter =
		"%(asctime)s %(levelname)s [%(name)s] %(message)s"

	dbLogger := logging.NewLogger("db")
	dbLogger.Parent = logger
	dbLogger.Level = logging.LEVEL_WARN

	debug := flag.Int("x", 0, "set debug modes, (default: 0)")
	backend := flag.String("backend", "",
		fmt.Sprintf("select backend [%s]", strings.Join(BACKENDS, "|")))
	setup := flag.Bool("setup", false,
		"perform database setup before operation")
	flag.Parse()

	if *debug > 0 {
		logger.Level = logging.LEVEL_DEBUG
	}
	if *debug >= 3 {
		dbLogger.Level = logging.LEVEL_DEBUG
	}

	// select database backend
	var InteractiveConfig func(KwArgs) (KwArgs, error)
	var InteractiveSetup func(KwArgs) error
	var NewHandler func(KwArgs) *db.Handler
	switch *backend {
	case "sqlite":
		DB_OPTIONS["database"] = fmt.Sprintf(
			"/tmp/%v.db", DB_OPTIONS["database"])
		InteractiveConfig = sqlite.InteractiveConfig
		InteractiveSetup = sqlite.InteractiveSetup
		NewHandler = sqlite.NewHandler
	case "mysql":
		DB_OPTIONS["port"] = 3306
		InteractiveConfig = mysql.InteractiveConfig
		InteractiveSetup = mysql.InteractiveSetup
		NewHandler = mysql.NewHandler
	case "pgsql":
		DB_OPTIONS["port"] = 5432
		DB_OPTIONS["username"] = "postgres"
		DB_OPTIONS["password"] = ""
		InteractiveConfig = pgsql.InteractiveConfig
		InteractiveSetup = pgsql.InteractiveSetup
		NewHandler = pgsql.NewHandler
	case "mssql":
		DB_OPTIONS["port"] = 1433
		DB_OPTIONS["username"] = "sa"
		DB_OPTIONS["password"] = "root@Root"
		InteractiveConfig = mssql.InteractiveConfig
		InteractiveSetup = mssql.InteractiveSetup
		NewHandler = mssql.NewHandler
	default:
		fmt.Print("\nError!! invalid database backend\n\n")
		return
	}

	fmt.Printf("\n* Using backend: %v\n", *backend)
	fmt.Println("\nConfig:")
	dbOpts, err := InteractiveConfig(DB_OPTIONS)
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			fmt.Print("\n--exit--\n\n")
		} else {
			fmt.Printf("Error: %v\n\n", err)
		}
		os.Exit(0)
	}

	fmt.Println("\nDB Options:")
	for _, k := range []string{"database",
		"host", "port", "username", "password", "extargs"} {
		fmt.Printf(" - %-9v: %v\n", k, dbOpts[k])
	}

	if *setup {
		fmt.Println("\nDB Setup:")
		if err := InteractiveSetup(dbOpts); err != nil {
			fmt.Printf("Error: %v\n\n", err)
			os.Exit(0)
		}
		fmt.Println("Done")
	}
	fmt.Println()

	// create database handler
	dbh := NewHandler(dbOpts)
	dbh.Logger = dbLogger

	runOperations(dbh)

	fmt.Println()
}
