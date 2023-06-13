package webapp

import (
	"net/http"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/exonlabs/go-utils/logging"
	loghandlers "github.com/exonlabs/go-utils/logging/handlers"
)

type WebServer struct {
	Name      string
	Options   KwArgs
	Log       *logging.Logger
	Reqlogger *logging.Logger
	Debug     uint8

	BasePath string
	// Views    []*WebView

	procPid int
	srvHnd  *http.Server
	srvMux  *http.ServeMux
}

func NewWebServer(name string, logger *logging.Logger,
	reqlogger *logging.Logger, Debug uint8, options KwArgs) *WebServer {
	return &WebServer{
		Name:      name,
		Options:   options,
		Log:       logger,
		Reqlogger: reqlogger,
		Debug:     Debug,
		// Views:     []*WebView{},
	}
}

func (ws *WebServer) Initialize() error {
	if ws.Log == nil {
		ws.Log = loghandlers.NewStdoutLogger(ws.Name)
		ws.Log.Formatter =
			"%(asctime)s - %(levelname)s [%(name)s] %(message)s"
	}

	ws.Log.Info("initializing")

	// TODO: set default req logger

	ws.srvMux = http.NewServeMux()
	if ws.BasePath != "" {
		ws.srvMux.Handle("/static/", http.StripPrefix(
			"/static/", http.FileServer(http.Dir(ws.BasePath))))
	}

	return nil
}

func (ws *WebServer) AddView(view *WebView) {

	// view.SetParams(ws)

	view.Parent = ws
	view.Log = ws.Log
	view.Debug = ws.Debug

	// if v, ok := view.(IViewInitialize); ok {
	// 	v.Initialize()
	// }

	for _, route := range view.Routes {
		ws.srvMux.HandleFunc(route, func(
			w http.ResponseWriter, r *http.Request) {

			ws.handleRequest(view, w, r)

			if ws.Reqlogger != nil {
				t := time.Now().Format("02/Jan/06 15:04:05")
				ws.Reqlogger.Info("%s - - [%s] \"%s %s %s\" %d -\n",
					strings.Split(r.RemoteAddr, ":")[0], t, r.Method,
					r.URL.String(), r.Proto, view.Env.Response.StatusCode)
			}
		})
	}
}

func (ws *WebServer) handleRequest(
	view *WebView, w http.ResponseWriter, r *http.Request) {

	env := &Env{
		Request:  &Request{r},
		Response: &Response{httpWriter: w},
	}

	// 	session, err := env.LoadSession()
	// 	if err != nil {
	// 		vh.parentSrv.Log.Error(err.Error())
	// 	}
	// 	if len(session) != 0 {
	// 		env.Session = session
	// 	} else {
	// 		env.Session = make(map[string]any)
	// 	}

	view.Env = env
	contents := view.handleRequest()

	if env.Response.StatusCode == 0 {
		env.Response.StatusCode = 200
	}

	// 	// generate session
	// 	if len(env.Session) != 0 {
	// 		err := env.setSession()
	// 		if err != nil {
	// 			vh.parentSrv.Log.Error(err.Error())
	// 		}
	// 	}

	// 	// check cookie
	// 	if len(env.cooikes) != 0 {
	// 		for _, cookie := range env.cooikes {
	// 			http.SetCookie(w, cookie)
	// 		}
	// 	}

	// 	// set http header
	// 	if len(env.Respone.Headers) != 0 {
	// 		for key, val := range env.Respone.Headers {
	// 			w.Header().Set(key, val)
	// 		}
	// 	}

	// set http status code
	w.WriteHeader(env.Response.StatusCode)

	_, err := w.Write([]byte(contents))
	if err != nil {
		ws.Log.Error(err.Error())
	}
}

func (ws *WebServer) Start(host string, port int) error {
	ws.procPid = syscall.Getpid()

	bind := host + ":" + strconv.Itoa(port)
	ws.srvHnd = &http.Server{
		Addr:    bind,
		Handler: ws.srvMux,
	}

	ws.Log.Info("running on http://%s", bind)
	return ws.srvHnd.ListenAndServe()
}

func (ws *WebServer) Stop() {
	if ws.procPid != 0 {
		ws.Log.Info("stop request")
		syscall.Kill(ws.procPid, syscall.SIGINT)
	}
}
