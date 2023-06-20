package web

import (
	"net/http"
	"strings"
)

type Request struct {
	*http.Request
}

// check xhr/ajax request type
func (r *Request) IsXhr() bool {
	return strings.Contains(
		r.Header.Get("X-Requested-With"), "XMLHttpRequest")
}

// check json or xhr/ajax request type
func (r *Request) IsJson() bool {
	return strings.Contains(r.Header.Get("Content-Type"), "json") ||
		strings.Contains(r.Header.Get("X-Requested-With"), "XMLHttpRequest")
}
