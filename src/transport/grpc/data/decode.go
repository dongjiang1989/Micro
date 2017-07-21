package Case

import (
	"golang.org/x/net/context"

	"transport/grpc/data/pb"
)

func decodeRequest(ctx context.Context, req interface{}) (interface{}, error) {
	r := req.(*pb.CaseRequest)
	return CaseRequest{A: r.A, B: r.B}, nil
}

func decodeResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	r := resp.(*pb.CaseResponse)
	return &CaseResponse{V: r.V, Ctx: ctx}, nil
}
