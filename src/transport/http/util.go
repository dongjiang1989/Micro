package http

import (
	"net/http"
)

// the Content-Type is set.
type Headerer interface {
	Headers() http.Header
}

// StatusCoder, the StatusCode will be used when encoding the error. By default,
// StatusInternalServerError (500) is used.
type StatusCoder interface {
	StatusCode() int
}
