package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/exonlabs/go-logging/pkg/xlog"
	"github.com/exonlabs/go-utils/pkg/webui"
)

var (
	HOST        = "0.0.0.0"
	PORT        = 8080
	WEBRES_PATH = filepath.Join(os.TempDir(), "web_resources")
	// menuBuffer = make(map[int]webui.MenuLink)
)

// prepare resources path
func prepareResources(localRes string) error {
	os.RemoveAll(WEBRES_PATH)
	err := os.MkdirAll(WEBRES_PATH, 0775)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return webui.InitResources(localRes)
}

func main() {
	logger := xlog.NewLogger("main")
	xlog.SetDefaultLogger(logger)
	xlog.SetDefaultFormatter("{time} {level} [{source}] {message}")

	debug := flag.Int("x", 0, "set debug modes, (default: 0)")
	// initWebui := flag.String(
	// 	"init-webui", "", "initialize webui resources in specified path")
	initRes := flag.String("init-resources", "", "initialize web resources")

	flag.Parse()

	if *debug > 0 {
		xlog.DefaultLogger().Level = int(*debug) * -10
	}

	logger.Info("***** starting *****")

	fmt.Printf("res: %s\n", *initRes)
	if *initRes != "" {
		if err := prepareResources(*initRes); err != nil {
			logger.Fatal("Failed to init resources, " + err.Error())
			os.Exit(1)
		}
	}
	// srv := web.NewServer("WebPortal", nil)
	// srv.DefaultContentType = "text/html; charset=utf-8"
	// // srv.StaticPath = STATIC_PATH

	// // srv.AddView(&IndexView{})
	// // srv.AddView(&HomeView{})
	// // srv.AddView(&ExitView{})

	// if err := srv.Start(HOST, PORT); err != nil {
	// 	logger.Fatal(err.Error())
	// 	os.Exit(1)
	// }
	// logger.Info("exit")

	// websrv := webapp.NewWebServer(
	// 	"WebSession", logger, reqlogger, uint8(*debug), nil)
	// websrv.BasePath = "static"
	// websrv.Initialize()

	// views := []*webapp.WebView{
	// 	IndexView, HomeView, NotifyView,
	// 	AlertsView, InputForm, DatagridView,
	// 	QueryBuilderView, LoaderView, LoginView,
	// }
	// for _, view := range views {
	// 	websrv.AddView(view)
	// }

	// if err := websrv.Start(HOST, PORT); err != nil {
	// 	logger.Error(err.Error())
	// }
}
