package smpp

import "errors"

var (
	ErrConnectionClosed = errors.New("smpp: connection closed")
)
