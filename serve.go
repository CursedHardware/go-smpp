package smpp

import (
	"context"
	"crypto/tls"
	"net"
)

type Handler interface{ Serve(*Session) }

type HandlerFunc func(*Session)

func (h HandlerFunc) Serve(session *Session) { h.Serve(session) }

func ServeTCP(address string, handler Handler, config *tls.Config) (err error) {
	var listener net.Listener
	if config == nil {
		listener, err = net.Listen("tcp", address)
	} else {
		listener, err = tls.Listen("tcp", address, config)
	}
	if err != nil {
		return
	}
	var parent net.Conn
	for {
		if parent, err = listener.Accept(); err != nil {
			return
		}
		go handler.Serve(NewSession(context.Background(), parent))
	}
}
