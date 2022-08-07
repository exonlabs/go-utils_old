package sqlite

import (
	"github.com/exonlabs/go-utils/pkg/db"
)


type Session struct {
	db.BaseSession
}


func NewSession(options map[string]any, logger any, debug uint8) *Session {
	this := Session{}
	this.Session = &this
	this.Backend = "sqlite"
	this.Options = options
	this.Debug = debug
	this.IsConnected = false

	// TODO: implement real logger
	// this.Logger = logger
	// switch {
	//     case this.Debug >= 5:
	//         this.Logger.SetLevel("DEBUG")
	//     case this.Debug >= 4:
	//         this.Logger.SetLevel("INFO")
	//     default:
	//         this.Logger.SetLevel("ERROR")
	// }

	return &this
}

func (this *Session) Connect() {
	// TODO

	this.IsConnected = true
}

func (this *Session) Close() {
	// TODO

	this.IsConnected = false
}

func (this *Session) DoExecute(query string, params ...any) (bool, string) {
	// TODO

	return true, ""
}

func (this *Session) DoExecuteScript(query_script string) (bool, string) {
	// TODO

	return true, ""
}

func (this *Session) FetchOne() *map[string]any {
	// TODO

	return nil
}

func (this *Session) FetchAll() *[]map[string]any {
	// TODO

	return nil
}

func (this *Session) RowCount() int32 {
	// TODO

	return -1
}

func (this *Session) LastRowId() int32 {
	// TODO

	return -1
}

func (this *Session) Commit() {
	// TODO
}

func (this *Session) RollBack() {
	// TODO
}
