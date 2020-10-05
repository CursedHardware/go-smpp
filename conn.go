package smpp

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"sync"
	"time"

	. "github.com/NiceLabs/go-smpp/pdu"
)

type Conn struct {
	parent       net.Conn
	ctx          context.Context
	cancel       context.CancelFunc
	sendQueue    chan interface{}
	receiveQueue chan interface{}
	pending      map[int32]func(interface{})
	mutex        sync.Mutex
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
		sendQueue:    make(chan interface{}),
		receiveQueue: make(chan interface{}),
		pending:      make(map[int32]func(interface{})),
		NextSequence: rand.Int31,
		ReadTimeout:  time.Minute * 15,
		WriteTimeout: time.Minute * 15,
	}
}

func (c *Conn) Watch() {
	defer c.cancel()
	go c.watchOutbound()
	c.watchInbound()
}

func (c *Conn) watchInbound() {
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
		if packet, err = ReadPDU(c.parent); err != nil {
			continue
		}
		if callback, ok := c.pending[ReadSequence(packet)]; ok {
			go callback(packet)
		} else {
			c.receiveQueue <- packet
		}
	}
}

func (c *Conn) watchOutbound() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case packet := <-c.sendQueue:
			if c.WriteTimeout > 0 {
				_ = c.parent.SetReadDeadline(time.Now().Add(c.WriteTimeout))
			}
			_, err := Marshal(c.parent, packet)
			if callback, ok := c.pending[ReadSequence(packet)]; ok {
				go callback(err)
			}
		}
	}
}

func (c *Conn) Submit(ctx context.Context, packet Responsable) (resp interface{}, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	sequence := c.NextSequence()
	WriteSequence(packet, sequence)
	returns := make(chan interface{}, 1)
	c.sendQueue <- packet
	c.pending[sequence] = func(resp interface{}) {
		if resp == nil {
			return
		}
		returns <- resp
		delete(c.pending, sequence)
	}
	select {
	case <-c.ctx.Done():
		delete(c.pending, sequence)
		err = c.ctx.Err()
	case <-ctx.Done():
		delete(c.pending, sequence)
		err = ctx.Err()
	case resp = <-returns:
		var ok bool
		if err, ok = resp.(error); ok {
			resp = nil
		}
	}
	return
}

func (c *Conn) Send(ctx context.Context, packet interface{}) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	sequence := ReadSequence(packet)
	if sequence == 0 || sequence < 0 {
		err = errors.New("smpp: sequence unset")
		return
	}
	returns := make(chan interface{}, 1)
	c.sendQueue <- packet
	c.pending[sequence] = func(resp interface{}) {
		returns <- resp
		delete(c.pending, sequence)
	}
	select {
	case <-c.ctx.Done():
		delete(c.pending, sequence)
		err = ctx.Err()
	case <-ctx.Done():
		delete(c.pending, sequence)
		err = ctx.Err()
	case resp := <-returns:
		err, _ = resp.(error)
	}
	return
}

func (c *Conn) EnquireLink(tick time.Duration, timeout time.Duration) {
	ticker := time.NewTicker(tick)
	defer ticker.Stop()
	send := func() (err error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err = c.Submit(ctx, new(EnquireLink))
		return
	}
	for {
		if err := send(); err == context.DeadlineExceeded {
			_ = c.Close()
			return
		}
		<-ticker.C
	}
}

func (c *Conn) Close() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	defer c.cancel()
	_, err = c.Submit(ctx, new(Unbind))
	if err == nil {
		close(c.sendQueue)
		close(c.receiveQueue)
		err = c.parent.Close()
	}
	return
}

func (c *Conn) PDU() <-chan interface{} { return c.receiveQueue }
