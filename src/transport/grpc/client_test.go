package grpc_test

import (
	"fmt"
	"net"
	"testing"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"

	"google.golang.org/grpc"

	Case "transport/grpc/data"

	Pb "transport/grpc/data/pb"
)

const (
	hostPort string = "localhost:46837"
)

func Test_GRPCClient(t *testing.T) {
	assert := assert.New(t)

	var (
		server  = grpc.NewServer()
		service = Case.NewService()
	)

	sc, err := net.Listen("tcp", hostPort)
	assert.Nil(err)

	defer server.GracefulStop()

	go func() {
		Pb.RegisterCaseServer(server, Case.NewBinding(service))
		_ = server.Serve(sc)
	}()

	cc, err := grpc.Dial(hostPort, grpc.WithInsecure())
	assert.Nil(err)

	client := Case.NewClient(cc)

	var (
		a   = "the answer to life the universe and everything"
		b   = int64(42)
		cID = "request-info"
		ctx = Case.SetCorrelationID(context.Background(), cID)
	)

	responseCTX, v, err := client.Case(ctx, a, b)
	assert.Nil(err)

	assert.Equal(fmt.Sprintf("%s = %d", a, b), v)

	assert.Equal(cID, Case.GetConsumedCorrelationID(responseCTX))
}
