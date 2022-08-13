package db

type IModel interface {
	TableName() string
	TableArgs() map[string]any
	TableColumns() [][]string
	TableConstraints() string
	DataAdapters(data map[string]any) map[string]any
	// DataConverters() map[string]func(any)any
	InitializeData(ISession)
	// UpgradeSchema()
	// DowngradeSchema()
}

type BaseModel struct {
	Table_Name        string
	Table_Args        map[string]any
	Table_Columns     [][]string
	Table_Constraints string
	Data_Adapters     map[string]func(any) any
}

func (this *BaseModel) TableName() string {
	return this.Table_Name
}

func (this *BaseModel) TableArgs() map[string]any {
	return this.Table_Args
}

func (this *BaseModel) TableColumns() [][]string {
	var res [][]string
	guid := false
	// check if table_column has "guid"
	for _, table_column := range this.Table_Columns {
		if table_column[0] == "guid" {
			guid = true
		}
	}
	// if not column append "guid"
	if !guid {
		res = [][]string{
			{"guid", "TEXT NOT NULL", "PRIMARY"},
		}
	}
	res = append(res, this.Table_Columns...)
	return res
}

func (this *BaseModel) TableConstraints() string {
	return this.Table_Constraints
}

func (this *BaseModel) DataAdapters(data map[string]any) map[string]any {
	for key, fn := range this.Data_Adapters {
		if val, ok := data[key]; ok {
			data[key] = fn(val)
		}
	}
	return data
}

func (this *BaseModel) InitializeData(ISession) {}
