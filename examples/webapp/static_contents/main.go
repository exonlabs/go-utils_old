package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/exonlabs/go-logging/pkg/xlog"
	"github.com/exonlabs/go-utils_old/pkg/web"
)

var (
	HOST        = "0.0.0.0"
	PORT        = 8000
	STATIC_PATH = filepath.Join(os.TempDir(), "web_static")
)

type IndexView struct{}

func (v *IndexView) Meta() *web.ViewMeta {
	return web.NewViewMeta("index", "/")
}

func (v *IndexView) DoGet(ctx *web.Context) *web.Response {
	content := `<a href="/home">Home Page</a>`
	return web.NewResponse(content)
}

type HomeView struct{}

func (v *HomeView) Meta() *web.ViewMeta {
	return web.NewViewMeta("home", "/home", "/home/")
}

func (v *HomeView) DoGet(ctx *web.Context) *web.Response {
	content := `<!doctype html>
<html>
<head>
<link rel="stylesheet" type="text/css" href="/static/css/test.css">
</head>
<body>
<p class="text">Home Page</p>
</body>
</html>`
	return web.NewResponse(content)
}

type ExitView struct{}

func (v *ExitView) Meta() *web.ViewMeta {
	return web.NewViewMeta("exit", "/exit")
}

func (v *ExitView) DoGet(ctx *web.Context) *web.Response {
	ctx.Server.Stop()
	return web.NewResponse("")
}

// create temp static contents dir and return its path
func prepareStatic() error {
	os.RemoveAll(STATIC_PATH)
	err := os.MkdirAll(STATIC_PATH, 0775)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// create sample css contents
	cssPath := filepath.Join(STATIC_PATH, "css")
	cssContent := "p.text {font-size:16px; font-weight:bold}"
	os.MkdirAll(cssPath, 0775)
	os.WriteFile(filepath.Join(cssPath, "test.css"), []byte(cssContent), 0664)

	return nil
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

	// set static folder
	if err := prepareStatic(); err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}
	defer os.RemoveAll(STATIC_PATH)

	srv := web.NewServer("WebPortal", nil)
	srv.DefaultContentType = "text/html; charset=utf-8"
	srv.StaticPath = STATIC_PATH

	srv.AddView(&IndexView{})
	srv.AddView(&HomeView{})
	srv.AddView(&ExitView{})

	if err := srv.Start(HOST, PORT); err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}
	logger.Info("exit")
}
