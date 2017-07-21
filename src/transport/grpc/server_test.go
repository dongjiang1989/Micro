package grpc_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	grpctransport "transport/grpc"
	Case "transport/grpc/data"

	log "github.com/cihub/seelog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	Pb "transport/grpc/data/pb"
)

func Test_ExampleGrpcServer(t *testing.T) {
	log.Flush()
	assert := assert.New(t)

	handler := grpctransport.NewServer(
		func(ctx context.Context, request interface{}) (interface{}, error) {
			req := request.(Case.CaseRequest)
			newCtx, v, err := Case.NewService().Case(ctx, req.A, req.B)
			return &Case.CaseResponse{
				V:   v,
				Ctx: newCtx,
			}, err
		},
		func(cxt context.Context, req interface{}) (interface{}, error) {
			log.Info("requerst: i am dongjiang's request!!!")
			r := req.(*Pb.CaseRequest)
			return Case.CaseRequest{A: r.A, B: r.B}, nil
		},
		func(cxt context.Context, resp interface{}) (interface{}, error) {
			log.Info("response: i am dongjiang's response!!!", resp)
			r := resp.(*Case.CaseResponse)
			return &Pb.CaseResponse{V: r.V}, nil
		},
	)
	server := grpc.NewServer()

	sc, err := net.Listen("tcp", "localhost:46838")
	assert.Nil(err)

	defer server.GracefulStop()

	go func() {
		Pb.RegisterCaseServer(server, Case.NewBindingHandler(handler))
		_ = server.Serve(sc)
	}()

	conn, err := grpc.Dial("localhost:46838", grpc.WithInsecure())
	assert.Nil(err)
	client := Case.NewClientHandler(grpctransport.NewClient(
		conn,
		"pb.Case",
		"Case",
		func(ctx context.Context, req interface{}) (interface{}, error) {
			log.Info("requerst: i am dongjiang's request", req)
			r := req.(Case.CaseRequest)
			return &Pb.CaseRequest{A: r.A, B: r.B}, nil
		},
		func(ctx context.Context, resp interface{}) (interface{}, error) {
			log.Info("response: i am dongjiang's response")
			r := resp.(*Pb.CaseResponse)
			return &Case.CaseResponse{V: r.V, Ctx: ctx}, nil
		},
		&Pb.CaseResponse{},
		grpctransport.ClientBefore(),
		grpctransport.ClientAfter(),
	))

	var (
		a   = "the answer to life the universe and everything"
		b   = int64(42)
		cID = "request-info"
		ctx = Case.SetCorrelationID(context.Background(), cID)
	)

	_, v, err := client.Case(ctx, a, b)
	assert.Nil(err)
	assert.Equal(fmt.Sprintf("%s = %d", a, b), v)

	assert.True(true)
}
