package main

import (
	"flag"

	"github.com/exonlabs/go-utils/logging"
	"github.com/exonlabs/go-utils/logging/handlers"
	"github.com/exonlabs/go-utils/webapp"
)

var (
	HOST = "0.0.0.0"
	PORT = 8000
)

var IndexView = &webapp.WebView{
	Name:   "index",
	Routes: []string{"/"},

	GET: func(v *webapp.WebView) string {
		v.Log.Info("test log msg from index")
		return v.Name + v.Env.Request.Host
	},
	POST: func(v *webapp.WebView) string {
		return ""
	},
}

// type IndexView struct {
// 	*webapp.BaseView
// }

// func (v *IndexView) GetRoutes() []string {
// 	return []string{"/"}
// }

// func (v *IndexView) DispatchRequest(env *webapp.Env) {
// 	env.Respone.Content = `<!doctype html>
// <html>
// <head>
// <link rel="stylesheet" type="text/css" href="/static/test.css">
// </head>
// <body>
// <p class="mybold">index page</p>
// </body>
// </html>
// `
// 	env.Respone.StatusCode = 200
// }

// type HomeView struct {
// 	*webapp.BaseView
// }

// func (v *HomeView) GetRoutes() []string {
// 	return []string{"/home/"}
// }

// func (v *HomeView) DispatchRequest(env *webapp.Env) {
// 	env.Respone.Content = `<!doctype html>
// <html>
// <head>
// <link rel="stylesheet" type="text/css" href="/static/test.css">
// </head>
// <body>
// <p class="mybold">home page</p>
// </body>
// </html>
// `
// 	env.Respone.StatusCode = 200
// }

// func (h *HomeView) Get(env *webapp.Env) {
// 	env.Respone.Content = h.GetUrlPrefix()
// 	env.Respone.StatusCode = 200
// }

// type ExitView struct {
// 	*webapp.BaseView
// }

// func (e *ExitView) Get(env *webapp.Env) {
// 	e.Parent.Stop()
// }

func main() {
	logger := handlers.NewStdoutLogger("main")

	reqlogger := handlers.NewStdoutLogger("reqlogger")
	reqlogger.Formatter = "%(message)s"

	debug := flag.Int("x", 0, "set debug modes, (default: 0)")
	flag.Parse()

	if *debug > 0 {
		logger.Level = logging.LEVEL_DEBUG
	}

	websrv := webapp.NewWebServer(
		"WebSession", logger, reqlogger, uint8(*debug), nil)
	websrv.BasePath = "/tmp/css"
	websrv.Initialize()

	websrv.AddView(IndexView)
	// websrv.AddView(&HomeView{})
	// websrv.AddView(&ExitView{})

	if err := websrv.Start(HOST, PORT); err != nil {
		logger.Error(err.Error())
	}
}
