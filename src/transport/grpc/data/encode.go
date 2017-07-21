package Case

import (
	"golang.org/x/net/context"

	"transport/grpc/data/pb"
)

func encodeRequest(ctx context.Context, req interface{}) (interface{}, error) {
	r := req.(CaseRequest)
	return &pb.CaseRequest{A: r.A, B: r.B}, nil
}

func encodeResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	r := resp.(*CaseResponse)
	return &pb.CaseResponse{V: r.V}, nil
}
