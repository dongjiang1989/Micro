package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	log "github.com/cihub/seelog"
)

type Response struct {
	Body   io.ReadCloser
	String string
}

func Test_HttpClient(t *testing.T) {
	log.Flush()
	assert := assert.New(t)

	var (
		testbody = "dongjiang‘s Micro"
		encode   = func(context.Context, *http.Request, interface{}) error { return nil }
		decode   = func(_ context.Context, r *http.Response) (interface{}, error) {
			buffer := make([]byte, len(testbody))
			r.Body.Read(buffer)
			return Response{r.Body, string(buffer)}, nil
		}
		headers        = make(chan string, 1)
		headerKey      = "who"
		headerVal      = "dongjiang"
		afterHeaderKey = "X-hereIsThere"
		afterHeaderVal = "dongjiang"
		afterVal       = ""
		afterFunc      = func(ctx context.Context, r *http.Response) context.Context {
			afterVal = r.Header.Get(afterHeaderKey)
			return ctx
		}
	)

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				headers <- r.Header.Get(headerKey)
				w.Header().Set(afterHeaderKey, afterHeaderVal)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(testbody))
			}))

	client := NewClient(
		"GET",
		mustParse(server.URL),
		encode,
		decode,
		ClientBefore(SetRequestHeader(headerKey, headerVal)),
		ClientAfter(afterFunc),
	)

	res, err := client.Endpoint()(context.Background(), struct{}{})
	assert.Nil(err)

	var have string
	select {
	case have = <-headers:
		log.Info("Do nothing! ", have)
	case <-time.After(time.Millisecond):
		assert.False(true, "not in this branch!")
		log.Critical("timeout waiting for %s", headerKey)
	}
	// Check that Request Header was successfully received
	assert.Equal(headerVal, have, "want is not have!")

	// Check that Response header set from server was received in SetClientAfter
	assert.Equal(afterVal, have, "want is not have!")

	// Check that the response was successfully decoded
	response, ok := res.(Response)
	assert.True(ok)

	want, have := testbody, response.String
	assert.Equal(want, have, "want is not have!")

	// Check that response body was closed
	b := make([]byte, 1)
	_, err = response.Body.Read(b)
	assert.NotNil(err)

	doNotWant, err := io.EOF, err
	assert.NotEqual(doNotWant, err)
}

func mustParse(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		log.Critical(err)
	}
	return u
}
