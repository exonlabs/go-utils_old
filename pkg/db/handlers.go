package db

type IDBHandler interface {
	CreateSession() ISession
	InitDatabase([]IModel) bool
}

type BaseDBHandler struct {
	Options map[string]any
	Logger  any
	Backend string
}

// return session handler object
func (this *BaseDBHandler) CreateSession() ISession {
	panic("NOT_IMPLEMENTED")
}

// create database tables and initialize data
func (this *BaseDBHandler) InitDatabase(models []IModel) bool {
	return false
}
