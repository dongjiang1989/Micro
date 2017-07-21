package grpc

import (
	"golang.org/x/net/context"

	"google.golang.org/grpc/metadata"
)

// ResponseFuncs are only executed in servers, after invoking the endpoint but prior to writing a response.
type ServerResponseFunc func(ctx context.Context, header *metadata.MD, trailer *metadata.MD) context.Context

// ClientResponseFuncs are only executed in clients, after a request has been made, but prior to it being decoded.
type ClientResponseFunc func(ctx context.Context, header metadata.MD, trailer metadata.MD) context.Context

// SetResponseHeader returns a ResponseFunc that sets the specified metadata
// key-value pair.
func SetResponseHeader(key, val string) ServerResponseFunc {
	return func(ctx context.Context, md *metadata.MD, _ *metadata.MD) context.Context {
		key, val := EncodeKeyValue(key, val)
		(*md)[key] = append((*md)[key], val)
		return ctx
	}
}
