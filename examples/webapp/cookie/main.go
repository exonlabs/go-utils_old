package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/exonlabs/go-utils/logging"
	"github.com/exonlabs/go-utils/logging/handlers"
	"github.com/exonlabs/go-utils/webapp"
)

const (
	HOST         = "0.0.0.0"
	PORT         = 8000
	cooikeEncKey = "password12345678"
)

type CookieView struct {
	*webapp.BaseView
}

func NewCookieView() *CookieView {
	return &CookieView{webapp.NewBaseView("index")}
}

func (c *CookieView) Get(env *webapp.ViewEnv) {
	env.SetCookie(&http.Cookie{Name: "msg", Value: "Hello"})
	env.SetCookie(&http.Cookie{Name: "foo", Value: "bar"})

	err := env.SetEncryptionCookie(cooikeEncKey, &http.Cookie{Name: "msg_encoded", Value: "Hello"})
	if err != nil {
		c.Log.Error(err.Error())
	}

	err = env.SetEncryptionCookie(cooikeEncKey, &http.Cookie{Name: "foo_encoded", Value: "bar"})
	if err != nil {
		c.Log.Error(err.Error())
	}

	env.Respone.Content = "set cookies"
}

type GetCookieView struct {
	*webapp.BaseView
}

func NewGetCookieView() *GetCookieView {
	return &GetCookieView{webapp.NewBaseView("getcookie")}
}

func (c *GetCookieView) Get(env *webapp.ViewEnv) {
	cookieNames := []string{"msg", "foo"}
	str := "get all cookies\n"
	for key := range cookieNames {
		cookie, err := env.GetCookie(cookieNames[key])
		if err != nil {
			c.Log.Error(err.Error())
			return
		}
		str += fmt.Sprintf("cookie clear name:%s, value:%s\n", cookie.Name, cookie.Value)
	}

	cookieNamesEnc := []string{"msg_encoded", "foo_encoded"}
	for key := range cookieNamesEnc {
		cookie, err := env.GetEncryptionCookie(cooikeEncKey, cookieNamesEnc[key])
		if err != nil {
			c.Log.Error(err.Error())
			return
		}
		str += fmt.Sprintf("cookie encoded name:%s, value:%s\n", cookie.Name, cookie.Value)
	}

	env.Respone.Content = str
}

type DelCookieView struct {
	*webapp.BaseView
}

func NewDelCookieView() *DelCookieView {
	return &DelCookieView{webapp.NewBaseView("delcookie")}
}

func (c *DelCookieView) Get(env *webapp.ViewEnv) {
	err := env.DelCookie("msg")
	if err != nil {
		c.Log.Error(err.Error())
		return
	}

	env.Respone.Content = "delete `msg` cookie"
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

	websrv := webapp.NewServer("WebSession", logger, reqlogger, uint8(*debug), nil)
	websrv.Initialize()

	views := []webapp.IView{
		NewCookieView(),
		NewGetCookieView(),
		NewDelCookieView(),
	}

	for _, v := range views {
		websrv.AddView(v)
	}

	if err := websrv.Start(HOST, PORT); err != nil {
		logger.Error(err.Error())
	}
}
