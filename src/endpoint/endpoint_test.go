package endpoint

import (
	"fmt"
	"testing"

	log "github.com/cihub/seelog"
	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
)

func Test_ExampleChain(t *testing.T) {
	log.Flush()
	assert := assert.New(t)

	e := Chain(
		annotate("1"),
		annotate("2"),
		annotate("3"),
		annotate("4"),
		annotate("5"),
		annotate("6"),
		annotate("7"),
		annotate("8"),
		annotate("9"),
		annotate("10"),
	)(myEndpoint)

	if _, err := e(ctx, req); err != nil {
		assert.NotNil(err)
	}
}

var (
	ctx = context.Background()
	req = struct {
		info      string
		dongjiang string
	}{

		info:      "aaa",
		dongjiang: "I am dongjiang",
	}
)

func annotate(s string) Middleware {
	return func(next Endpoint) Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {

			fmt.Println(s, "pre")
			log.Info(s, "pre")

			defer func() {
				log.Info(s, "post")
				fmt.Println(s, "post")
			}()

			return next(ctx, request)
		}
	}
}

func myEndpoint(cxt context.Context, rqt interface{}) (interface{}, error) {
	log.Info("dongjiang endpoint!")
	return Nop(cxt, rqt)
}
