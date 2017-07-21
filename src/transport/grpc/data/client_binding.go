package Case

import (
	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"endpoint"

	grpctransport "transport/grpc"

	"transport/grpc/data/pb"
)

type clientBinding struct {
	e endpoint.Endpoint
}

func (c *clientBinding) Case(ctx context.Context, a string, b int64) (context.Context, string, error) {
	response, err := c.e(ctx, CaseRequest{A: a, B: b})
	if err != nil {
		return nil, "", err
	}
	r := response.(*CaseResponse)
	return r.Ctx, r.V, nil
}

func NewClientHandler(handler *(grpctransport.Client)) Service {
	return &clientBinding{e: handler.Endpoint()}
}

func NewClient(conn *grpc.ClientConn) Service {
	return &clientBinding{
		e: grpctransport.NewClient(
			conn,
			"pb.Case",
			"Case",
			encodeRequest,
			decodeResponse,
			&pb.CaseResponse{},
			grpctransport.ClientBefore(
				injectCorrelationID,
			),
			grpctransport.ClientBefore(
				displayClientRequestHeaders,
			),
			grpctransport.ClientAfter(
				displayClientResponseHeaders,
				displayClientResponseTrailers,
			),
			grpctransport.ClientAfter(
				extractConsumedCorrelationID,
			),
		).Endpoint(),
	}
}
