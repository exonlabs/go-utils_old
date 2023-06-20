package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/exonlabs/go-logging/pkg/xlog"
	"github.com/exonlabs/go-utils/pkg/web"
)

var (
	HOST       = "0.0.0.0"
	PORT       = 8000
	SESS_STORE = filepath.Join(os.TempDir(), "web_sessions")
)

type IndexView struct{}

func (v *IndexView) Meta() *web.ViewMeta {
	return web.NewViewMeta("index", "/")
}

func (v *IndexView) DoGet(ctx *web.Context) *web.Response {
	content := "Session Data:\n"
	for key, val := range ctx.Session.Buffer() {
		content += fmt.Sprintf("%s: %v\n", key, val)
	}
	return web.NewResponse(content)
}

type SessionAddView struct{}

func (v *SessionAddView) Meta() *web.ViewMeta {
	return web.NewViewMeta("add", "/add")
}

func (v *SessionAddView) DoGet(ctx *web.Context) *web.Response {
	oldVal := 0.0
	if v, ok := ctx.Session.Get("int_counter"); ok {
		oldVal = v.(float64)
	}
	ctx.Session.Set("int_counter", oldVal+1)
	ctx.Session.Set("float_value", 12.3456)
	ctx.Session.Set("string_value", "some string value")
	ctx.Session.Set("bool_value", true)
	return web.RedirectResponse("/")
}

type SessionSubView struct{}

func (v *SessionSubView) Meta() *web.ViewMeta {
	return web.NewViewMeta("sub", "/sub")
}

func (v *SessionSubView) DoGet(ctx *web.Context) *web.Response {
	oldVal := 0.0
	if v, ok := ctx.Session.Get("int_counter"); ok {
		oldVal = v.(float64)
	}
	ctx.Session.Set("int_counter", oldVal-1)
	return web.RedirectResponse("/")
}

type SessionDelView struct{}

func (v *SessionDelView) Meta() *web.ViewMeta {
	return web.NewViewMeta("del", "/del")
}

func (v *SessionDelView) DoGet(ctx *web.Context) *web.Response {
	ctx.Session.Del("int_counter")
	ctx.Session.Set("float_value", 0.0)
	ctx.Session.Set("string_value", "")
	ctx.Session.Set("bool_value", false)
	return web.RedirectResponse("/")
}

type SessionPurgeView struct{}

func (v *SessionPurgeView) Meta() *web.ViewMeta {
	return web.NewViewMeta("purge", "/purge")
}

func (v *SessionPurgeView) DoGet(ctx *web.Context) *web.Response {
	ctx.Session.Purge()
	return web.RedirectResponse("/")
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
	filestore := flag.Bool("file-store", false, "use file-base session store")
	flag.Parse()

	if *debug > 0 {
		xlog.DefaultLogger().Level = int(*debug) * -10
	}

	logger.Info("***** starting *****")

	srv := web.NewServer("WebPortal", nil)
	srv.DefaultContentType = "text/plain; charset=utf-8"
	if *filestore {
		logger.Info("-- using file-based session store --")
		if _, err := os.Stat(SESS_STORE); err != nil {
			if err := os.MkdirAll(SESS_STORE, 0775); err != nil {
				panic(err)
			}
		}
		srv.SessionFactory = web.NewSessionFileFactory(map[string]any{
			// "secret_key": "0123456789abcdef",
			"store_path": SESS_STORE,
		})
	} else {
		srv.SessionFactory = web.NewSessionCookieFactory(map[string]any{
			"secret_key": "0123456789abcdef",
		})
	}

	srv.AddView(&IndexView{})
	srv.AddView(&SessionAddView{})
	srv.AddView(&SessionSubView{})
	srv.AddView(&SessionDelView{})
	srv.AddView(&SessionPurgeView{})
	srv.AddView(&ExitView{})

	if err := srv.Start(HOST, PORT); err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}
	logger.Info("exit")
}
