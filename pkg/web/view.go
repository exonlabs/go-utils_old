package web

import (
	"net/http"
)

type ViewMeta struct {
	Name   string
	Routes []string
}

func NewViewMeta(name string, routes ...string) *ViewMeta {
	return &ViewMeta{
		Name:   name,
		Routes: routes,
	}
}

type View interface {
	Meta() *ViewMeta
}

type ViewDispatchRequest interface {
	DispatchRequest(*Context) *Response
}
type ViewBeforeRequest interface {
	BeforeRequest(*Context) *Response
}
type ViewAfterRequest interface {
	AfterRequest(*Context, *Response) *Response
}
type ViewHandleRequest interface {
	HandleRequest(*Context) *Response
}
type ViewDoHead interface {
	DoHead(*Context) *Response
}
type ViewDoGet interface {
	DoGet(*Context) *Response
}
type ViewDoPost interface {
	DoPost(*Context) *Response
}
type ViewDoPut interface {
	DoPut(*Context) *Response
}
type ViewDoPatch interface {
	DoPatch(*Context) *Response
}
type ViewDoDelete interface {
	DoDelete(*Context) *Response
}

func dispatchRequest(ctx *Context, view View) *Response {
	// delegate execution to view dispatch method if exist
	if v, ok := view.(ViewDispatchRequest); ok {
		return v.DispatchRequest(ctx)
	}

	var resp *Response

	// run before request method if exist and return if there is
	// generated response after call
	if v, ok := view.(ViewBeforeRequest); ok {
		resp = v.BeforeRequest(ctx)
		if resp != nil {
			return resp
		}
	}

	// run handle request method if exist
	if v, ok := view.(ViewHandleRequest); ok {
		resp = v.HandleRequest(ctx)
	} else {
		// run method matching the http method
		switch ctx.Request.Method {
		case http.MethodPost:
			if v, ok := view.(ViewDoPost); ok {
				resp = v.DoPost(ctx)
			}
		case http.MethodGet:
			if v, ok := view.(ViewDoGet); ok {
				resp = v.DoGet(ctx)
			}
		case http.MethodHead:
			if v, ok := view.(ViewDoHead); ok {
				resp = v.DoHead(ctx)
			} else if v, ok := view.(ViewDoGet); ok {
				// run GET as alternative to head
				resp = v.DoGet(ctx)
			}
		case http.MethodPut:
			if v, ok := view.(ViewDoPut); ok {
				resp = v.DoPut(ctx)
			}
		case http.MethodPatch:
			if v, ok := view.(ViewDoPatch); ok {
				resp = v.DoPatch(ctx)
			}
		case http.MethodDelete:
			if v, ok := view.(ViewDoDelete); ok {
				resp = v.DoDelete(ctx)
			}
		default:
			return ErrorResponse(
				"Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}

	// run after request method if exist
	if v, ok := view.(ViewAfterRequest); ok {
		return v.AfterRequest(ctx, resp)
	}

	return resp
}
