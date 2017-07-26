package client

import (
	"time"

	"golang.org/x/net/context"
)

//-----------------------------------------------
type Option func(*Options)

//-----------------------------------------------

type Client interface {
	Init(...Option)
	Options() Options
}
