package http

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"

	"endpoint"

	log "github.com/cihub/seelog"
	"github.com/stretchr/testify/assert"
)

func Test_ServerBadDecode(t *testing.T) {

	assert := assert.New(t)

	handler := NewServer(
		func(context.Context, interface{}) (interface{}, error) { return struct{}{}, nil },
		func(context.Context, *http.Request) (interface{}, error) {
			return struct{}{}, errors.New("duangduangduang")
		},
		func(context.Context, http.ResponseWriter, interface{}) error { return nil },
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	log.Info("Url: ", server.URL)

	resp, err := http.Get(server.URL)
	assert.Nil(err)

	want, have := http.StatusInternalServerError, resp.StatusCode
	assert.Equal(want, have)
}

func Test_ServerBadEndpoint(t *testing.T) {
	assert := assert.New(t)

	handler := NewServer(
		func(context.Context, interface{}) (interface{}, error) {
			return struct{}{}, errors.New("duangduangduang")
		},
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		func(context.Context, http.ResponseWriter, interface{}) error { return nil },
	)
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Info("Url: ", server.URL)
	resp, err := http.Get(server.URL)
	assert.Nil(err)

	want, have := http.StatusInternalServerError, resp.StatusCode
	assert.Equal(want, have)
}

func Test_ServerBadEncode(t *testing.T) {
	assert := assert.New(t)

	handler := NewServer(
		func(context.Context, interface{}) (interface{}, error) { return struct{}{}, nil },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		func(context.Context, http.ResponseWriter, interface{}) error { return errors.New("dang") },
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	log.Info("Url: ", server.URL)

	resp, err := http.Get(server.URL)
	assert.Nil(err)
	want, have := http.StatusInternalServerError, resp.StatusCode
	assert.Equal(want, have)
}

func Test_ServerErrorEncoder(t *testing.T) {
	assert := assert.New(t)

	errTeapot := errors.New("dongjiang")

	code := func(err error) int {
		if err == errTeapot {
			return http.StatusTeapot
		}
		return http.StatusInternalServerError
	}

	handler := NewServer(
		func(context.Context, interface{}) (interface{}, error) { return struct{}{}, errTeapot },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		func(context.Context, http.ResponseWriter, interface{}) error { return nil },
		ServerErrorEncoder(func(_ context.Context, err error, w http.ResponseWriter) { w.WriteHeader(code(err)) }),
	)
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Info(server.URL)
	resp, err := http.Get(server.URL)
	assert.Nil(err)

	want, have := http.StatusTeapot, resp.StatusCode
	assert.Equal(want, have)
}

func Test_ServerHappyPath(t *testing.T) {
	assert := assert.New(t)

	step, response := MockServer(t)
	step()
	resp := <-response
	defer resp.Body.Close()
	_, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	want, have := http.StatusOK, resp.StatusCode
	assert.Equal(want, have)
}

func Test_MultipleServerBefore(t *testing.T) {
	assert := assert.New(t)

	var (
		headerKey    = "X-Henlo-Lizer"
		headerVal    = "Helllo you stinky lizard"
		statusCode   = http.StatusTeapot
		responseBody = "go eat a fly ugly\n"
		done         = make(chan struct{})
	)
	handler := NewServer(
		endpoint.Nop,
		func(context.Context, *http.Request) (interface{}, error) {
			return struct{}{}, nil
		},
		func(_ context.Context, w http.ResponseWriter, _ interface{}) error {
			w.Header().Set(headerKey, headerVal)
			w.WriteHeader(statusCode)
			w.Write([]byte(responseBody))
			return nil
		},
		ServerBefore(func(ctx context.Context, r *http.Request) context.Context {
			ctx = context.WithValue(ctx, "one", 1)
			return ctx
		}),
		ServerBefore(func(ctx context.Context, r *http.Request) context.Context {
			if _, ok := ctx.Value("one").(int); !ok {
				t.Error("Value was not set properly when multiple ServerBefores are used")
			}
			close(done)
			return ctx
		}),
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	go http.Get(server.URL)

	select {
	case <-done:
	case <-time.After(time.Second):
		assert.True(false)
	}
}

func Test_MultipleServerAfter(t *testing.T) {
	assert := assert.New(t)

	var (
		headerKey    = "X-Henlo-Lizer"
		headerVal    = "Helllo you stinky lizard"
		statusCode   = http.StatusTeapot
		responseBody = "go eat a fly ugly\n"
		done         = make(chan struct{})
	)

	handler := NewServer(
		endpoint.Nop,
		func(context.Context, *http.Request) (interface{}, error) {
			return struct{}{}, nil
		},
		func(_ context.Context, w http.ResponseWriter, _ interface{}) error {
			w.Header().Set(headerKey, headerVal)
			w.WriteHeader(statusCode)
			w.Write([]byte(responseBody))
			return nil
		},
		ServerAfter(func(ctx context.Context, w http.ResponseWriter) context.Context {
			ctx = context.WithValue(ctx, "one", 1)
			return ctx
		}),
		ServerAfter(func(ctx context.Context, w http.ResponseWriter) context.Context {
			_, ok := ctx.Value("one").(int)
			assert.True(ok, "Value was not set properly when multiple ServerAfters are used")
			close(done)
			return ctx
		}),
	)

	server := httptest.NewServer(handler)
	defer server.Close()
	go http.Get(server.URL)

	select {
	case <-done:
	case <-time.After(time.Second):
		assert.True(false)
	}
}

func Test_ServerFinalizer(t *testing.T) {
	assert := assert.New(t)

	var (
		headerKey    = "X-Henlo-Lizer"
		headerVal    = "Helllo you stinky lizard"
		statusCode   = http.StatusTeapot
		responseBody = "go eat a fly ugly\n"
		done         = make(chan struct{})
	)
	handler := NewServer(
		endpoint.Nop,
		func(context.Context, *http.Request) (interface{}, error) {
			return struct{}{}, nil
		},
		func(_ context.Context, w http.ResponseWriter, _ interface{}) error {
			w.Header().Set(headerKey, headerVal)
			w.WriteHeader(statusCode)
			w.Write([]byte(responseBody))
			return nil
		},
		ServerFinalizer(func(ctx context.Context, code int, _ *http.Request) {
			want, have := statusCode, code
			assert.Equal(want, have)

			responseHeader := ctx.Value(ContextKeyResponseHeaders).(http.Header)
			want1, have1 := headerVal, responseHeader.Get(headerKey)
			assert.Equal(want1, have1)

			responseSize := ctx.Value(ContextKeyResponseSize).(int64)
			want2, have2 := int64(len(responseBody)), responseSize
			assert.Equal(want2, have2)

			close(done)
		}),
	)

	server := httptest.NewServer(handler)
	defer server.Close()
	go http.Get(server.URL)

	select {
	case <-done:
	case <-time.After(time.Second):
		assert.True(false, "timeout waiting for finalizer")
	}
}

type enhancedResponse struct {
	Foo string `json:"foo"`
}

func (e enhancedResponse) StatusCode() int      { return http.StatusPaymentRequired }
func (e enhancedResponse) Headers() http.Header { return http.Header{"X-Edward": []string{"Snowden"}} }

func Test_EncodeJSONResponse(t *testing.T) {
	assert := assert.New(t)

	handler := NewServer(
		func(context.Context, interface{}) (interface{}, error) { return enhancedResponse{Foo: "bar"}, nil },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		EncodeJSONResponse,
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	assert.Nil(err)

	assert.Equal(http.StatusPaymentRequired, resp.StatusCode)

	assert.Equal("Snowden", resp.Header.Get("X-Edward"))

	buf, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(`{"foo":"bar"}`, strings.TrimSpace(string(buf)))
}

type noContentResponse struct{}

func (e noContentResponse) StatusCode() int { return http.StatusNoContent }

func Test_EncodeNoContent(t *testing.T) {
	assert := assert.New(t)

	handler := NewServer(
		func(context.Context, interface{}) (interface{}, error) { return noContentResponse{}, nil },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		EncodeJSONResponse,
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, resp.StatusCode)

	buf, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(0, len(buf))
}

type enhancedError struct{}

func (e enhancedError) Error() string                { return "enhanced error" }
func (e enhancedError) StatusCode() int              { return http.StatusTeapot }
func (e enhancedError) MarshalJSON() ([]byte, error) { return []byte(`{"err":"enhanced"}`), nil }
func (e enhancedError) Headers() http.Header         { return http.Header{"X-Enhanced": []string{"1"}} }

func Test_EnhancedError(t *testing.T) {
	assert := assert.New(t)

	handler := NewServer(
		func(context.Context, interface{}) (interface{}, error) { return nil, enhancedError{} },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		func(_ context.Context, w http.ResponseWriter, _ interface{}) error { return nil },
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	assert.Nil(err)

	defer resp.Body.Close()
	assert.Equal(http.StatusTeapot, resp.StatusCode)
	assert.Equal("1", resp.Header.Get("X-Enhanced"))

	buf, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(`{"err":"enhanced"}`, strings.TrimSpace(string(buf)))
}

func MockServer(t *testing.T) (step func(), resp <-chan *http.Response) {
	var (
		stepch   = make(chan bool)
		endpoint = func(context.Context, interface{}) (interface{}, error) { <-stepch; return struct{}{}, nil }
		response = make(chan *http.Response)

		handler = NewServer(
			endpoint,
			func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
			func(context.Context, http.ResponseWriter, interface{}) error { return nil },
			ServerBefore(func(ctx context.Context, r *http.Request) context.Context { return ctx }),
			ServerAfter(func(ctx context.Context, w http.ResponseWriter) context.Context { return ctx }),
		)
	)

	go func() {
		server := httptest.NewServer(handler)
		defer server.Close()

		resp, err := http.Get(server.URL)

		if err != nil {
			t.Error(err)
			return
		}

		response <- resp
	}()

	return func() { stepch <- true }, response
}
