package web

import (
	"github.com/exonlabs/go-logging/pkg/xlog"
)

type Context struct {
	// server handler
	Server *Server
	// log handler for request
	Logger *xlog.Logger
	// request handler
	Request *Request

	// session store
	Session SessionStore
}
