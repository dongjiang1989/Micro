package http

import (
	"net/http"

	"golang.org/x/net/context"
)

type EncodeRequestFunc func(context.Context, *http.Request, interface{}) error

type EncodeResponseFunc func(context.Context, http.ResponseWriter, interface{}) error
