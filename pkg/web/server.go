package web

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/exonlabs/go-logging/pkg/xlog"
)

type Server struct {
	Name      string
	Logger    *xlog.Logger
	Reqlogger *xlog.Logger

	srvHnd *http.Server
	srvMux *http.ServeMux

	// server session factory
	SessionFactory SessionFactory

	// path for serving static contents
	StaticPath string
	// max request content length
	MaxContentLength int
	// default response content type
	DefaultContentType string
}

func NewServer(name string, opts map[string]any) *Server {
	srv := &Server{
		Name:             name,
		srvMux:           http.NewServeMux(),
		MaxContentLength: 10485760, // 10 MiB: 10*1024*1024
	}

	if v, ok := opts["static_path"]; ok {
		srv.StaticPath = v.(string)
	}
	if v, ok := opts["max_content_length"]; ok {
		srv.MaxContentLength = v.(int)
	}
	if v, ok := opts["default_content_type"]; ok {
		srv.DefaultContentType = v.(string)
	}

	return srv
}

func (s *Server) AddView(view View) {
	meta := view.Meta()
	for _, route := range meta.Routes {
		s.srvMux.HandleFunc(
			route, func(w http.ResponseWriter, r *http.Request) {
				s.handleRequest(meta.Name, view, r, w)
			})
	}
}

func (s *Server) handleRequest(
	name string, view View, r *http.Request, w http.ResponseWriter) {

	// recover panic
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("InternalServerError"))
			s.Logger.Error(fmt.Sprintf("%s", err))
			s.Logger.Trace(string(debug.Stack()))
		}
	}()

	var sessStore SessionStore
	if s.SessionFactory != nil {
		sessStore = s.SessionFactory.Create(r, w)
		if err := sessStore.Load(); err != nil {
			s.Logger.Warn("failed loading session, " + err.Error())
		}
	}

	ctx := &Context{
		Server:  s,
		Logger:  s.Logger.CreateChild(name),
		Request: &Request{Request: r},
		Session: sessStore,
	}

	// handle request
	resp := dispatchRequest(ctx, view)
	if resp == nil {
		resp = ErrorResponse("Bad Request", http.StatusBadRequest)
	}

	// check and set response headers
	if resp.Headers != nil && len(resp.Headers) > 0 {
		for k, v := range resp.Headers {
			w.Header().Set(k, v)
		}
	}
	// set default Content-Type
	if v := w.Header().Get("Content-Type"); v == "" {
		w.Header().Set("Content-Type", s.DefaultContentType)
	}

	// check and add response cookies
	if resp.Cookies != nil && len(resp.Cookies) > 0 {
		for _, c := range resp.Cookies {
			w.Header().Add("Set-Cookie", c.String())
		}
	}

	// check and save session
	if sessStore != nil {
		if err := sessStore.Save(); err != nil {
			s.Logger.Warn("failed saving session, " + err.Error())
		}
	}

	if resp.RedirectUrl != "" {
		http.Redirect(w, r, resp.RedirectUrl, resp.StatusCode)
		s.reqLog(r, resp.StatusCode, 0, resp.RedirectUrl)
	} else {
		if resp.StatusCode == 0 {
			resp.StatusCode = http.StatusOK
		}

		// write response header with status code
		w.WriteHeader(resp.StatusCode)

		// write contents
		size, err := w.Write(resp.Content)
		if err != nil {
			s.Logger.Error(err.Error())
		}
		s.reqLog(r, resp.StatusCode, size, "")
	}
}

func (s *Server) reqLog(r *http.Request, code int, size int, refer string) {
	if s.Reqlogger != nil {
		ts := time.Now().Format("02/Jan/06 15:04:05.000")
		msg := fmt.Sprintf("%s - - [%s] \"%s %s %s\" %d %d",
			strings.Split(r.RemoteAddr, ":")[0], ts, r.Method,
			r.URL.String(), r.Proto, code, size)
		if refer != "" {
			msg += fmt.Sprintf(" \"%s\"", refer)
		}
		s.Reqlogger.Info(msg)
	}
}

func (s *Server) Start(host string, port int) error {
	// adjust logging
	if s.Logger == nil {
		s.Logger = xlog.DefaultLogger()
	}
	if s.Reqlogger == nil {
		s.Reqlogger = xlog.NewLogger("reqlogger")
	}
	s.Reqlogger.SetFormatter("{message}")

	// adjust static routing mux
	if s.StaticPath != "" {
		s.srvMux.Handle("/static/", http.StripPrefix(
			"/static/", http.FileServer(http.Dir(s.StaticPath))))
	}

	bind := host + ":" + strconv.Itoa(port)
	s.Logger.Info("running on http://%s", bind)

	s.srvHnd = &http.Server{
		Addr:    bind,
		Handler: s.srvMux,
	}
	err := s.srvHnd.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	if s.srvHnd != nil {
		s.Logger.Info("stop request")
		err := s.srvHnd.Close()
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}
