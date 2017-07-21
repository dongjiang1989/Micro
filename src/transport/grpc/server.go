package grpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"endpoint"

	log "github.com/cihub/seelog"
)

// Handler which should be called from the gRPC binding of the service implementation.
type Handler interface {
	ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) // grpc service
}

// Server wraps an endpoint and implements grpc.Handler.
type Server struct {
	e      endpoint.Endpoint
	dec    DecodeRequestFunc
	enc    EncodeResponseFunc
	before []ServerRequestFunc
	after  []ServerResponseFunc
}

// NewServer constructs a new server, which implements wraps the provided endpoint and implements the Handler interface. Consumers should write
// bindings that adapt the concrete gRPC methods from their compiled protobuf
// definitions to individual handlers. Request and response objects are from the
// caller business domain, not gRPC request and reply types.
func NewServer(
	e endpoint.Endpoint,
	dec DecodeRequestFunc,
	enc EncodeResponseFunc,
	options ...ServerOption,
) *Server {
	s := &Server{
		e:   e,
		dec: dec,
		enc: enc,
	}
	for _, option := range options {
		option(s)
	}
	return s
}

// ServerOption sets an optional parameter for servers.
type ServerOption func(*Server)

// ServerBefore functions are executed on the HTTP request object before the
// request is decoded.
func ServerBefore(before ...ServerRequestFunc) ServerOption {
	return func(s *Server) { s.before = append(s.before, before...) }
}

// ServerAfter functions are executed on the HTTP response writer after the
// endpoint is invoked, but before anything is written to the client.
func ServerAfter(after ...ServerResponseFunc) ServerOption {
	return func(s *Server) { s.after = append(s.after, after...) }
}

// ServeGRPC implements grpc.Handler.
func (s Server) ServeGRPC(ctx context.Context, req interface{}) (context.Context, interface{}, error) {
	// Retrieve gRPC metadata.
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	for _, f := range s.before {
		ctx = f(ctx, md)
	}

	request, err := s.dec(ctx, req)
	if err != nil {
		log.Error("Err: ", err)
		return ctx, nil, err
	}

	response, err := s.e(ctx, request)
	if err != nil {
		log.Error("Err: ", err)
		return ctx, nil, err
	}

	var mdHeader, mdTrailer metadata.MD
	for _, f := range s.after {
		ctx = f(ctx, &mdHeader, &mdTrailer)
	}

	grpcResp, err := s.enc(ctx, response)
	if err != nil {
		log.Error("Err: ", err)
		return ctx, nil, err
	}

	if len(mdHeader) > 0 {
		if err = grpc.SendHeader(ctx, mdHeader); err != nil {
			log.Error("Err: ", err)
			return ctx, nil, err
		}
	}

	if len(mdTrailer) > 0 {
		if err = grpc.SetTrailer(ctx, mdTrailer); err != nil {
			log.Error("Err: ", err)
			return ctx, nil, err
		}
	}

	return ctx, grpcResp, nil
}
