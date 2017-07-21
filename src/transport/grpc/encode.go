package grpc

import (
	"golang.org/x/net/context"
)

// EncodeRequestFunc encodes the passed request object into the gRPC request
// object. It's designed to be used in gRPC clients, for client-side endpoints.
type EncodeRequestFunc func(context.Context, interface{}) (request interface{}, err error)

// EncodeResponseFunc encodes the passed response object to the gRPC response
// message. It's designed to be used in gRPC servers, for server-side endpoints.
type EncodeResponseFunc func(context.Context, interface{}) (response interface{}, err error)
