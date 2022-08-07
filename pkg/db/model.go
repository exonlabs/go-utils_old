package db


type IModel interface {
	TableName() string
	TableArgs() map[string]any
	TableColumns() [][]string
	TableConstraints() []string
	// DataAdapters() map[string]func(any)any
	// DataConverters() map[string]func(any)any
	// InitializeData()
	// UpgradeSchema()
	// DowngradeSchema()
}

type BaseModel struct {
	IModel

	Table_Name string
	Table_Args map[string]any
	Table_Columns [][]string
	Table_Constraints []string
	// Data_Adapters map[string]func(any)any
	// Data_Converters map[string]func(any)any
	// Initialize_Data
	// Upgrade_Schema
	// Downgrade_Schema
}


func (this *BaseModel) TableName() string {
	return this.Table_Name
}

func (this *BaseModel) TableArgs() map[string]any {
	return this.Table_Args
}

func (this *BaseModel) TableColumns() [][]string {
	return this.Table_Columns
}

func (this *BaseModel) TableConstraints() []string {
	return this.Table_Constraints
}
