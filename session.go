package smpp

import (
	"context"
	"io"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/M2MGateway/go-smpp/pdu"
)

type Session struct {
	parent       net.Conn
	receiveQueue chan any
	pending      *sync.Map
	NextSequence func() int32
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewSession(ctx context.Context, parent net.Conn) (session *Session) {
	random := rand.New(rand.NewSource(time.Now().Unix()))
	session = &Session{
		parent:       parent,
		receiveQueue: make(chan any),
		pending:      new(sync.Map),
		NextSequence: random.Int31,
		ReadTimeout:  time.Minute * 15,
		WriteTimeout: time.Minute * 15,
	}
	go session.watch(ctx)
	return
}

//goland:noinspection SpellCheckingInspection
func (c *Session) watch(ctx context.Context) {
	var err error
	var packet any
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if c.ReadTimeout > 0 {
			_ = c.parent.SetReadDeadline(time.Now().Add(c.ReadTimeout))
		}
		if packet, err = pdu.Unmarshal(c.parent); err == io.EOF {
			return
		}
		if packet == nil {
			continue
		}
		if status, ok := err.(pdu.CommandStatus); ok {
			_ = c.Send(&pdu.GenericNACK{
				Header: pdu.Header{CommandStatus: status, Sequence: pdu.ReadSequence(packet)},
				Tags:   pdu.Tags{0xFFFF: []byte(err.Error())},
			})
			continue
		}
		if callback, ok := c.pending.Load(pdu.ReadSequence(packet)); ok {
			callback.(func(any))(packet)
		} else {
			c.receiveQueue <- packet
		}
	}
}

func (c *Session) Submit(ctx context.Context, packet pdu.Responsable) (resp any, err error) {
	sequence := c.NextSequence()
	pdu.WriteSequence(packet, sequence)
	if err = c.Send(packet); err != nil {
		return
	}
	returns := make(chan any, 1)
	c.pending.Store(sequence, func(resp any) { returns <- resp })
	select {
	case <-ctx.Done():
		err = ErrConnectionClosed
	case resp = <-returns:
	}
	c.pending.Delete(sequence)
	return
}

func (c *Session) Send(packet any) (err error) {
	sequence := pdu.ReadSequence(packet)
	if sequence == 0 || sequence < 0 {
		err = pdu.ErrInvalidSequence
		return
	}
	if c.WriteTimeout > 0 {
		err = c.parent.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
	}
	if err == nil {
		_, err = pdu.Marshal(c.parent, packet)
	}
	if err == io.EOF {
		err = ErrConnectionClosed
	}
	return
}

func (c *Session) EnquireLink(ctx context.Context, tick time.Duration, timeout time.Duration) (err error) {
	ticker := time.NewTicker(tick)
	defer ticker.Stop()
	for {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		if _, err = c.Submit(ctx, new(pdu.EnquireLink)); err != nil {
			ticker.Stop()
			err = c.Close(ctx)
		}
		cancel()
		<-ticker.C
	}
}

func (c *Session) Close(ctx context.Context) (err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	_, err = c.Submit(ctx, new(pdu.Unbind))
	if err != nil {
		return
	}
	close(c.receiveQueue)
	return c.parent.Close()
}

func (c *Session) PDU() <-chan any {
	return c.receiveQueue
}
