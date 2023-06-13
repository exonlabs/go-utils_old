package main

import (
	"flag"
	"fmt"

	"github.com/exonlabs/go-utils/logging"
	"github.com/exonlabs/go-utils/logging/handlers"
	"github.com/exonlabs/go-utils/webapp"
)

var (
	HOST          = "0.0.0.0"
	PORT          = 8000
	SessionEncKey = "PasswordPassword"
)

type SessionView struct {
	*webapp.BaseView
}

func NewSessionView() *SessionView {
	return &SessionView{webapp.NewBaseView("index")}
}

func (v *SessionView) Get(env *webapp.ViewEnv) {
	env.Session["username"] = "Mohammad"
	env.Session["userrole"] = "Admin"
	env.Respone.Content = "add sessions"
}

type GetSessionView struct {
	*webapp.BaseView
}

func NewGetSessionView() *GetSessionView {
	return &GetSessionView{webapp.NewBaseView("getsession")}
}

func (s *GetSessionView) Get(env *webapp.ViewEnv) {
	env.Session["foo"] = "bar"
	delete(env.Session, "userrole")
	str := "get all sessiones\n"
	sessions, err := env.LoadSession()
	if err != nil {
		s.Log.Error(err.Error())
		return
	}
	for key, value := range sessions {
		str += fmt.Sprintf("session key:%s, value:%s\n", key, value)
	}

	env.Respone.Content = str
}

func main() {
	logger := handlers.NewStdoutLogger("main")

	reqlogger := handlers.NewStdoutLogger("reqlogger")
	reqlogger.Formatter = "%(message)s"

	debug := flag.Int("x", 0, "set debug modes, (default: 0)")
	flag.Parse()

	if *debug > 0 {
		logger.Level = logging.LEVEL_DEBUG
	}

	ops := map[string]any{"secret_key": SessionEncKey}
	websrv := webapp.NewServer("WebSession", logger, reqlogger, uint8(*debug), ops)
	websrv.Initialize()

	views := []webapp.IView{
		NewSessionView(),
		NewGetSessionView(),
	}

	for _, v := range views {
		websrv.AddView(v)
	}

	if err := websrv.Start(HOST, PORT); err != nil {
		logger.Error(err.Error())
	}
}
