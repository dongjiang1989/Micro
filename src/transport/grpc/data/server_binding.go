package Case

import (
	"fmt"

	"golang.org/x/net/context"

	"endpoint"

	grpctransport "transport/grpc"

	"transport/grpc/data/pb"
)

type service struct{}

func (service) Case(ctx context.Context, a string, b int64) (context.Context, string, error) {
	return nil, fmt.Sprintf("%s = %d", a, b), nil
}

func NewService() Service {
	return service{}
}

func makeCaseEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CaseRequest)
		newCtx, v, err := svc.Case(ctx, req.A, req.B)
		return &CaseResponse{
			V:   v,
			Ctx: newCtx,
		}, err
	}
}

type serverBinding struct {
	handler grpctransport.Handler
}

func (b *serverBinding) Case(ctx context.Context, req *pb.CaseRequest) (*pb.CaseResponse, error) {
	_, response, err := b.handler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return response.(*pb.CaseResponse), nil
}

func NewBindingHandler(handler grpctransport.Handler) *serverBinding {
	return &serverBinding{
		handler: handler,
	}
}

func NewBinding(svc Service) *serverBinding {
	return &serverBinding{
		handler: grpctransport.NewServer(
			makeCaseEndpoint(svc),
			decodeRequest,
			encodeResponse,
			grpctransport.ServerBefore(
				extractCorrelationID,
			),
			grpctransport.ServerBefore(
				displayServerRequestHeaders,
			),
			grpctransport.ServerAfter(
				injectResponseHeader,
				injectResponseTrailer,
				injectConsumedCorrelationID,
			),
			grpctransport.ServerAfter(
				displayServerResponseHeaders,
				displayServerResponseTrailers,
			),
		),
	}
}
