package client

import (
	"codec"
	"time"

	"golang.org/x/net/context"
)

type Options struct {
	// Used to select codec
	ContentType string

	// Plugged interfaces
	Broker    broker.Broker
	Codecs    map[string]codec.NewCodec
	Registry  registry.Registry
	Selector  selector.Selector
	Transport transport.Transport

	// Connection Pool
	PoolSize int
	PoolTTL  time.Duration

	// Middleware for client
	Wrappers []Wrapper

	// Default Call Options
	CallOptions CallOptions

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

func (p *Options) InitNew() *Options {
	return &Options{}
}

func (p *Options) Hash() string {
	return abc
}
