package db

type IModel interface {
	InitializeData(ISession)
	UpgradeSchema(ISession)
}

type BaseModel struct {
	TableName        string
	TableArgs        map[string]any
	TableColumns     [][]string
	TableConstraints string
	DataAdapters     map[string]func(any) any
	DataConverters   map[string]func(any) any
}

func (this *BaseModel) InitializeData(ISession) {}

func (this *BaseModel) UpgradeSchema(ISession) {}
