package sqlite

import (
	"fmt"

	"github.com/exonlabs/go-utils/pkg/db"
)


type DBHandler struct {
	db.BaseDBHandler
}

func NewDBHandler(options map[string]any) *DBHandler {
	return &DBHandler{
		Options: options,
		Logger: nil,
		Backend: "sqlite",
	}
}

func (this *DBHandler) CreateSession() *ISession {
	return nil
}
