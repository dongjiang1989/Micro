package http

import (
	"net/http"

	"golang.org/x/net/context"
)

type DecodeRequestFunc func(context.Context, *http.Request) (request interface{}, err error)

type DecodeResponseFunc func(context.Context, *http.Response) (response interface{}, err error)
