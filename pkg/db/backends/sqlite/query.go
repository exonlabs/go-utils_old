package sqlite

import (
	"strings"

	"github.com/exonlabs/go-utils/pkg/db"
)

type Query struct {
	db.BaseQuery
}

func NewQuery(dbs db.ISession, model db.IModel) *Query {
	var this Query
	this.IQuery = &this
	this.DBs = dbs
	this.Model = model
	return &this
}

func (this *Query) CreateTable() {
	q := "CREATE TABLE IF NOT EXISTS " + this.Model.TableName() + "("
	for _, tableCoulmn := range this.Model.TableColumns() {
		for _, v := range tableCoulmn {
			q += v
			q += " "
		}
		q += ","
	}
	if len(this.Model.TableConstraints()) > 0 {
		q += this.Model.TableConstraints()
	}
	q = strings.TrimSpace(q)
	q = strings.TrimSuffix(q, ",")
	q += ")"
	this.DBs.Connect()
	this.DBs.Execute(q, false)
	this.DBs.Close()
}
