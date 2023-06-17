package main

import (
	"flag"

	"github.com/exonlabs/go-logging/pkg/xlog"
	"github.com/exonlabs/go-utils/pkg/webapp"
	"github.com/exonlabs/go-utils/pkg/webui"
)

var (
	HOST       = "0.0.0.0"
	PORT       = 8080
	menuBuffer = make(map[int]webui.MenuLink)
)

func main() {
	logger := xlog.NewLogger("main")
	logger.Formatter =
		"{time} {level} [{source}] {message}"

	reqlogger := xlog.NewLogger("reqlogger")
	reqlogger.Formatter = "{message}"

	debug := flag.Int("x", 0, "set debug modes, (default: 0)")
	flag.Parse()

	if *debug > 0 {
		logger.Level = xlog.LevelDebug
	}

	websrv := webapp.NewWebServer(
		"WebSession", logger, reqlogger, uint8(*debug), nil)
	websrv.BasePath = "static"
	websrv.Initialize()

	views := []*webapp.WebView{
		IndexView, HomeView, NotifyView,
		AlertsView, InputForm, DatagridView,
		QueryBuilderView, LoaderView, LoginView,
	}
	for _, view := range views {
		websrv.AddView(view)
	}

	if err := websrv.Start(HOST, PORT); err != nil {
		logger.Error(err.Error())
	}
}
