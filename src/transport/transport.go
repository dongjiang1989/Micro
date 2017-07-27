package transport

import (
	"time"
)

type Message struct {
	Header map[string]string
	Body   []byte
}

type Socket interface {
	Recv(*Message) error // receive message
	Send(*Message) error // send message
	Close() error
}

type Client interface {
	Socket
}

type Listener interface {
	Addr() string
	Close() error
	Accept(func(Socket)) error
}

type Option func(*Options)

type DialOption func(*DialOptions)

type ListenOption func(*ListenOptions)

// Transport is an interface which is used for communication between services.
// like: {HTTP, RabbitMQ, NATS, ...}
type Transport interface {
	Dial(addr string, opts ...DialOption) (Client, error)
	Listen(addr string, opts ...ListenOption) (Listener, error)
	String() string
}
