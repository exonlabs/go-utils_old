package web

import "net/http"

type Response struct {
	Headers     map[string]string
	Cookies     map[string]*http.Cookie
	Content     []byte
	RedirectUrl string
	StatusCode  int
}

func NewResponse(content string) *Response {
	return &Response{
		Content:    []byte(content),
		StatusCode: 200,
	}
}

func ErrorResponse(errmsg string, code int) *Response {
	return &Response{
		Content:    []byte(errmsg),
		StatusCode: code,
	}
}

func RedirectResponse(url string) *Response {
	return &Response{
		RedirectUrl: url,
		StatusCode:  http.StatusFound,
	}
}

// add new contents to end of current content
func (r *Response) Append(content string) {
	r.Content = append(r.Content, []byte(content)...)
}

// add new contents at begining of current content
func (r *Response) Prepend(content string) {
	r.Content = append([]byte(content), r.Content...)
}

// add/edit response header field in buffer
func (r *Response) SetHeader(key string, val string) {
	if key != "" && val != "" {
		if r.Headers == nil {
			r.Headers = make(map[string]string)
		}
		r.Headers[key] = val
	}
}

// delete response header field from buffer
func (r *Response) DelHeader(key string) {
	if r.Headers != nil {
		delete(r.Headers, key)
	}
}

// delete all response header fields from buffer
func (r *Response) FlushHeaders() {
	r.Headers = nil
}

// set/create new response cookie in buffer
func (r *Response) SetCookie(cookie *http.Cookie) {
	if cookie != nil {
		if r.Cookies == nil {
			r.Cookies = make(map[string]*http.Cookie)
		}
		r.Cookies[cookie.Name] = cookie
	}
}

// set/edit cookie to be deleted, adjust MaxAge=-1
func (r *Response) DelCookie(name string) {
	r.SetCookie(&http.Cookie{
		Name:   name,
		MaxAge: -1,
	})
}

// delete all response cookie fields from buffer
func (r *Response) FlushCookies() {
	r.Cookies = nil
}
