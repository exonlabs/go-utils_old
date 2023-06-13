package webapp

import (
	"net/http"
)

type Env struct {
	// option  map[string]any
	Request  *Request
	Response *Response
	// cooikes map[string]*http.Cookie
	// Session map[string]any
}

type Request struct {
	*http.Request
}

type Response struct {
	Headers    map[string]string
	StatusCode int
	httpWriter http.ResponseWriter
}

// func (resp *Response) SetHeader(key, val string) {
// 	if len(resp.Headers) == 0 {
// 		resp.Headers = make(map[string]string)
// 	}
// 	resp.Headers[key] = val
// }
