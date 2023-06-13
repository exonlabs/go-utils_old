package webapp

import (
	"net/http"
	"strings"

	"github.com/exonlabs/go-utils/logging"
)

type WebView struct {
	Name   string
	Parent *WebServer
	Log    *logging.Logger
	Debug  uint8

	Routes []string

	DispatchRequest func(*WebView) string
	// BeforeRequest func(*WebView) string
	// AfterRequest func(*WebView, string) string
	GET  func(*WebView) string
	POST func(*WebView) string
	// PUT func(*WebView, *Env) string
	// DELETE func(*WebView, *Env) string

	Env *Env
}

func (v *WebView) handleRequest() string {
	var contents string

	if v.DispatchRequest != nil {
		contents = v.DispatchRequest(v)
	} else if v.GET != nil {
		contents = v.GET(v)
	} else {
		contents = "Method Not Allowed"
		v.Env.Response.StatusCode = http.StatusMethodNotAllowed
	}

	return contents
}

// check xhr/ajax request type
func (v *WebView) IsXhrRequest() bool {
	h := v.Env.Request.Header
	return strings.Contains(h.Get("X-Requested-With"), "XMLHttpRequest")
}

// json or xhr/ajax request type
func (v *WebView) IsJsRequest() bool {
	h := v.Env.Request.Header
	return strings.Contains(h.Get("Content-Type"), "json") ||
		strings.Contains(h.Get("X-Requested-With"), "XMLHttpRequest")
}

func (v *WebView) Redirect(url string) {
	http.Redirect(
		v.Env.Response.httpWriter, v.Env.Request.Request,
		url, http.StatusFound)
}
