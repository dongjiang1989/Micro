package transport

import (
	"crypto/tls"
	"time"

	"transport/codec"

	"golang.org/x/net/context"
)

type Options struct {
	Addrs     []string
	Codec     codec.Codec
	Secure    bool
	TLSConfig *tls.Config
	// Timeout sets the timeout for Send/Recv
	Timeout time.Duration
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type DialOptions struct {
	Stream  bool
	Timeout time.Duration

	// Currently set in global options
	Context context.Context
}

type ListenOptions struct {
	// TODO: add tls options when listening
	// Currently set in global options

	Context context.Context
}

func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

func Codec(c codec.Codec) Option {
	return func(o *Options) {
		o.Codec = c
	}
}

func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

func WithStream() DialOption {
	return func(o *DialOptions) {
		o.Stream = true
	}
}

func WithTimeout(d time.Duration) DialOption {
	return func(o *DialOptions) {
		o.Timeout = d
	}
}
