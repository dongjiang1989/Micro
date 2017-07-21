package grpc

import (
	"golang.org/x/net/context"

	"google.golang.org/grpc/metadata"
)

const (
	binHdrSuffix = "-bin"
)

// ClientRequestFuncs are executed after creating the request but prior to sending the gRPC request to the server.
type ClientRequestFunc func(context.Context, *metadata.MD) context.Context

// ServerRequestFuncs are executed prior to invoking the endpoint.
type ServerRequestFunc func(context.Context, metadata.MD) context.Context

// SetRequestHeader returns a ClientRequestFunc that sets the specified metadata
// key-value pair.
func SetRequestHeader(key, val string) ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		key, val := EncodeKeyValue(key, val)
		(*md)[key] = append((*md)[key], val)
		return ctx
	}
}
