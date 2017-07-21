package http

import (
	"net/http"

	"golang.org/x/net/context"
)

//ServerResponseFuncs are only executed in servers, after invoking the endpoint but prior to writing a response.
type ServerResponseFunc func(context.Context, http.ResponseWriter) context.Context

//  ClientResponseFuncs are only executed in clients, after a request has been made, but prior to it being decoded.
type ClientResponseFunc func(context.Context, *http.Response) context.Context

func SetContentType(contentType string) ServerResponseFunc {
	return SetResponseHeader("Content-Type", contentType)
}

// SetResponseHeader returns a ServerResponseFunc that sets the given header.
func SetResponseHeader(key, val string) ServerResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter) context.Context {
		w.Header().Set(key, val)
		return ctx
	}
}
