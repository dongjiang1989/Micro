package grpc

import (
	"golang.org/x/net/context"
)

// DecodeRequestFunc extracts a user-domain request object from a gRPC request.
type DecodeRequestFunc func(context.Context, interface{}) (request interface{}, err error)

// DecodeResponseFunc extracts a user-domain response object from a gRPC
type DecodeResponseFunc func(context.Context, interface{}) (response interface{}, err error)
