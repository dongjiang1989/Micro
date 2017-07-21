package httproxy

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func Test_ServerHappyPathSingleServer(t *testing.T) {
	assert := assert.New(t)

	originServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("i am dongjiang!"))
	}))
	defer originServer.Close()

	originURL, err := url.Parse(originServer.URL)
	assert.Nil(err)

	handler := NewServer(
		originURL,
	)

	proxyServer := httptest.NewServer(handler)
	defer proxyServer.Close()

	resp, err := http.Get(proxyServer.URL)
	assert.Nil(err)

	assert.Equal(http.StatusOK, resp.StatusCode)

	responseBody, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("i am dongjiang!", string(responseBody))
}

func Test_ServerHappyPathSingleServerWithServerOptions(t *testing.T) {
	assert := assert.New(t)

	const (
		headerKey = "X-xxx"
		headerVal = "go-proxy"
	)

	originServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(headerVal, r.Header.Get(headerKey))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("i am dongjiang!"))
	}))
	defer originServer.Close()
	originURL, err := url.Parse(originServer.URL)
	assert.Nil(err)

	handler := NewServer(
		originURL,
		ServerBefore(func(ctx context.Context, r *http.Request) context.Context {
			r.Header.Add(headerKey, headerVal)
			return ctx
		}),
	)

	proxyServer := httptest.NewServer(handler)
	defer proxyServer.Close()

	resp, err := http.Get(proxyServer.URL)
	assert.Nil(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	responseBody, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("i am dongjiang!", string(responseBody))
}

func Test_ServerOriginServerNotFoundResponse(t *testing.T) {
	assert := assert.New(t)

	originServer := httptest.NewServer(http.NotFoundHandler())
	defer originServer.Close()

	originURL, err := url.Parse(originServer.URL)
	assert.Nil(err)

	handler := NewServer(
		originURL,
	)
	proxyServer := httptest.NewServer(handler)
	defer proxyServer.Close()

	resp, err := http.Get(proxyServer.URL)
	assert.Equal(http.StatusNotFound, resp.StatusCode)
}

func Test_ServerOriginServerUnreachable(t *testing.T) {
	assert := assert.New(t)

	// create a server, then promptly shut it down
	originServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	originURL, _ := url.Parse(originServer.URL)
	originServer.Close()

	handler := NewServer(
		originURL,
	)
	proxyServer := httptest.NewServer(handler)
	defer proxyServer.Close()

	resp, err := http.Get(proxyServer.URL)
	assert.Nil(err)

	switch resp.StatusCode {
	case http.StatusBadGateway: // go1.7 and beyond
		break
	case http.StatusInternalServerError: // to go1.7
		break
	default:
		assert.NotEqual(http.StatusBadGateway, resp.StatusCode)
		assert.NotEqual(http.StatusInternalServerError, resp.StatusCode)
	}
}

func Test_MultipleServerBefore(t *testing.T) {
	assert := assert.New(t)

	const (
		headerKey = "X-dj"
		headerVal = "go-proxy"
	)

	originServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(headerVal, r.Header.Get(headerKey))
		assert.Equal("bb", r.Header.Get("aa"))

		w.Header().Add(headerKey, headerVal)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("i'm dongjiang!"))
	}))
	defer originServer.Close()
	originURL, err := url.Parse(originServer.URL)
	assert.Nil(err)

	handler := NewServer(
		originURL,
		ServerBefore(func(ctx context.Context, r *http.Request) context.Context {
			r.Header.Add(headerKey, headerVal)
			r.Header.Add("aa", "bb")
			return ctx
		}),
		ServerBefore(func(ctx context.Context, r *http.Request) context.Context {
			return ctx
		}),
	)
	proxyServer := httptest.NewServer(handler)
	defer proxyServer.Close()

	resp, err := http.Get(proxyServer.URL)
	assert.Nil(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(headerVal, resp.Header.Get(headerKey))

	responseBody, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("i'm dongjiang!", string(responseBody))
}
