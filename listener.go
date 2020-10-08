package smpp

import (
	"context"
	"net"
)

func WatchListener(listener net.Listener, on func(*Conn)) (err error) {
	defer listener.Close()
	var parent net.Conn
	for {
		if parent, err = listener.Accept(); err != nil {
			return
		}
		go on(NewConn(context.Background(), parent))
	}
}
