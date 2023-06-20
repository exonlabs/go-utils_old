package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

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
	content := "All Cookies:\n"
	for _, c := range ctx.Request.Cookies() {
		content += fmt.Sprintf("%s: %v\n", c.Name, c.Value)
	}
	return web.NewResponse(content)
}

type CookieAddView struct{}

func (v *CookieAddView) Meta() *web.ViewMeta {
	return web.NewViewMeta("add", "/add")
}

func (v *CookieAddView) DoGet(ctx *web.Context) *web.Response {
	oldVal := 0
	if c, _ := ctx.Request.Cookie("counter"); c != nil {
		oldVal, _ = strconv.Atoi(c.Value)
	}
	resp := web.RedirectResponse("/")
	resp.SetCookie(&http.Cookie{
		Name:  "counter",
		Value: strconv.Itoa(oldVal + 1),
	})
	return resp
}

type CookieSubView struct{}

func (v *CookieSubView) Meta() *web.ViewMeta {
	return web.NewViewMeta("sub", "/sub")
}

func (v *CookieSubView) DoGet(ctx *web.Context) *web.Response {
	oldVal := 0
	if c, _ := ctx.Request.Cookie("counter"); c != nil {
		oldVal, _ = strconv.Atoi(c.Value)
	}
	resp := web.RedirectResponse("/")
	resp.SetCookie(&http.Cookie{
		Name:  "counter",
		Value: strconv.Itoa(oldVal - 1),
	})
	return resp
}

type CookieDelView struct{}

func (v *CookieDelView) Meta() *web.ViewMeta {
	return web.NewViewMeta("del", "/del")
}

func (v *CookieDelView) DoGet(ctx *web.Context) *web.Response {
	resp := web.RedirectResponse("/")
	resp.DelCookie("counter")
	return resp
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
	srv.DefaultContentType = "text/plain; charset=utf-8"

	srv.AddView(&IndexView{})
	srv.AddView(&CookieAddView{})
	srv.AddView(&CookieSubView{})
	srv.AddView(&CookieDelView{})
	srv.AddView(&ExitView{})

	if err := srv.Start(HOST, PORT); err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}
	logger.Info("exit")
}
