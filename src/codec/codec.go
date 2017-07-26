package codec

import (
	"io"
)

type MessageType int

const (
	Error       MessageType = iota // 0
	Request                        // 1
	Response                       // 2
	Publication                    // 3
)

// Takes in a connection/buffer and returns a new Codec
type NewCodec func(io.ReadWriteCloser) Codec

// Codec encodes/decodes various types of messages
type Codec interface {
	ReadHeader(*Message, MessageType) error
	ReadBody(interface{}) error // ReadHeader and ReadBody are called in pairs
	Write(*Message, interface{}) error
	Close() error
	String() string
}

// Message represents detailed information about the communication.
type Message struct {
	Id     uint64
	Type   MessageType
	Target string
	Method string
	Error  string
	Header map[string]string
}
