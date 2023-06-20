package main

import (
	"flag"
	"os"

	"github.com/exonlabs/go-logging/pkg/xlog"
	"github.com/exonlabs/go-utils/pkg/web"
)

var (
	HOST = "0.0.0.0"
	PORT = 8000
)

type IndexView struct{}

func (v *IndexView) Meta() *web.ViewMeta {
	return web.NewViewMeta("index", "/")
}

func (v *IndexView) DoGet(ctx *web.Context) *web.Response {
	ctx.Logger.Info("log msg from index GET")
	return web.NewResponse("host: " + ctx.Request.Host)
}

type HomeView struct{}

func (v *HomeView) Meta() *web.ViewMeta {
	return web.NewViewMeta("home", "/home", "/home/")
}

func (v *HomeView) DoGet(ctx *web.Context) *web.Response {
	ctx.Logger.Info("log msg from home GET")
	return web.NewResponse("this is home page")
}

type RedirectView struct{}

func (v *RedirectView) Meta() *web.ViewMeta {
	return web.NewViewMeta("redirect", "/redirect")
}

func (v *RedirectView) DoGet(ctx *web.Context) *web.Response {
	ctx.Logger.Info("log msg from redirect GET")
	return web.RedirectResponse("/home")
}

type ExitView struct{}

func (v *ExitView) Meta() *web.ViewMeta {
	return web.NewViewMeta("exit", "/exit")
}

func (v *ExitView) DoGet(ctx *web.Context) *web.Response {
	ctx.Server.Stop()
	return web.NewResponse("")
}

func main() {
	logger := xlog.NewLogger("main")
	xlog.SetDefaultLogger(logger)
	xlog.SetDefaultFormatter("{time} {level} [{source}] {message}")

	debug := flag.Int("x", 0, "set debug modes, (default: 0)")
	flag.Parse()

	if *debug > 0 {
		xlog.DefaultLogger().Level = int(*debug) * -10
	}

	logger.Info("***** starting *****")

	srv := web.NewServer("WebPortal", nil)

	srv.AddView(&IndexView{})
	srv.AddView(&HomeView{})
	srv.AddView(&RedirectView{})
	srv.AddView(&ExitView{})

	if err := srv.Start(HOST, PORT); err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}
	logger.Info("exit")
}
