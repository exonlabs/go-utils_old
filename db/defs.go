package db

import "database/sql"

const SQL_PLACEHOLDER = "$?"

// generic types
type KwArgs = map[string]any

// types defining table
type TableName = string
type TableColumns = [][]string
type TableConstraints = [][]string
type TableMeta = struct {
	Options     KwArgs
	Columns     TableColumns
	Constraints TableConstraints
}

// types for table data, where colums are mapped as keys and
// table conversion functions for data read and write
type BaseModel = struct{}
type ModelData = map[string]any
type ModelDataMappers = map[string](func(any) any)

// backend engine interface
type IEngine interface {
	GetBackendName() string
	FormatSqlStmt(string) string
	Connect(KwArgs) (*sql.DB, error)
	GenTableSchema(TableName, TableMeta) ([]string, error)
	ListRetryErrors() []string
}

// models interface
type IModel interface {
	GetTableName() TableName
	GetTableMeta() TableMeta
}
type IModelDefaultOrders interface {
	GetDefaultOrders() []string
}
type IModelDataAdapters interface {
	GetDataAdapters() ModelDataMappers
}
type IModelDataConverters interface {
	GetDataConverters() ModelDataMappers
}
type IModelInitializeData interface {
	InitializeData(*Session, TableName) error
}
type IModelUpgradeTableSchema interface {
	UpgradeTableSchema(*Session, TableName) error
}
