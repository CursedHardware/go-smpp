package smpp

import (
	"context"
	"io"
	"math/rand"
	"net"
	"time"

	. "github.com/NiceLabs/go-smpp/pdu"
)

type Conn struct {
	parent       net.Conn
	ctx          context.Context
	cancel       context.CancelFunc
	receiveQueue chan interface{}
	pending      map[int32]func(interface{})
	NextSequence func() int32
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewConn(ctx context.Context, parent net.Conn) *Conn {
	ctx, cancel := context.WithCancel(ctx)
	return &Conn{
		parent:       parent,
		ctx:          ctx,
		cancel:       cancel,
		receiveQueue: make(chan interface{}),
		pending:      make(map[int32]func(interface{})),
		NextSequence: rand.Int31,
		ReadTimeout:  time.Minute * 15,
		WriteTimeout: time.Minute * 15,
	}
}

func (c *Conn) Watch() {
	defer c.cancel()
	var err error
	var packet interface{}
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}
		if c.ReadTimeout > 0 {
			_ = c.parent.SetReadDeadline(time.Now().Add(c.ReadTimeout))
		}
		if packet, err = ReadPDU(c.parent); err == io.EOF {
			return
		} else if err != nil {
			continue
		} else if callback, ok := c.pending[ReadSequence(packet)]; ok {
			callback(packet)
		} else {
			c.receiveQueue <- packet
		}
	}
}

func (c *Conn) Submit(ctx context.Context, packet Responsable) (resp interface{}, err error) {
	sequence := c.NextSequence()
	WriteSequence(packet, sequence)
	if err = c.Send(packet); err != nil {
		return
	}
	returns := make(chan interface{}, 1)
	c.pending[sequence] = func(resp interface{}) { returns <- resp }
	defer delete(c.pending, sequence)
	select {
	case <-c.ctx.Done():
		err = ErrConnectionClosed
	case <-ctx.Done():
		err = ctx.Err()
	case resp = <-returns:
	}
	return
}

func (c *Conn) Send(packet interface{}) (err error) {
	sequence := ReadSequence(packet)
	if sequence == 0 || sequence < 0 {
		err = ErrInvalidSequence
		return
	}
	if c.WriteTimeout > 0 {
		err = c.parent.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
	}
	if err == nil {
		_, err = Marshal(c.parent, packet)
	}
	if err == io.EOF {
		err = ErrConnectionClosed
	}
	return
}

func (c *Conn) EnquireLink(tick time.Duration, timeout time.Duration) {
	ticker := time.NewTicker(tick)
	defer ticker.Stop()
	sendEnquireLink := func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := c.Submit(ctx, new(EnquireLink))
		if err == context.DeadlineExceeded || err == context.Canceled {
			ticker.Stop()
			_ = c.Close()
		}
	}
	for {
		sendEnquireLink()
		<-ticker.C
	}
}

func (c *Conn) Close() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	defer c.cancel()
	if _, err = c.Submit(ctx, new(Unbind)); err == nil {
		close(c.receiveQueue)
		err = c.parent.Close()
	}
	return
}

func (c *Conn) PDU() <-chan interface{} {
	return c.receiveQueue
}
