package smpp

import (
	"context"
	"io"
	"math/rand"
	"net"
	"time"

	. "github.com/M2MGateway/go-smpp/pdu"
)

type Session struct {
	parent       net.Conn
	receiveQueue chan Packet
	pending      map[int32]func(Packet)
	NextSequence func() int32
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Dial(ctx context.Context, address string) (conn *Session, err error) {
	var dialer net.Dialer
	parent, err := dialer.DialContext(ctx, "tcp", address)
	if err == nil {
		conn = NewSession(parent)
	}
	return
}

func DialTimeout(address string, timeout time.Duration) (*Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return Dial(ctx, address)
}

func NewSession(parent net.Conn) (conn *Session) {
	conn = &Session{
		parent:       parent,
		receiveQueue: make(chan Packet),
		pending:      make(map[int32]func(Packet)),
		NextSequence: rand.Int31,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}
	go conn.watch()
	return
}

//goland:noinspection SpellCheckingInspection
func (s *Session) watch() {
	var err error
	var packet Packet
	for {
		if s.ReadTimeout > 0 {
			_ = s.parent.SetReadDeadline(time.Now().Add(s.ReadTimeout))
		}
		packet, err = ReadPDU(s.parent)
		if err == io.EOF {
			break
		}
		if packet == nil {
			continue
		}
		sequence := ReadSequence(packet)
		if status, ok := err.(CommandStatus); err != nil {
			if !ok {
				status = ErrUnknownError
			}
			_ = s.Send(&GenericNACK{
				Header: Header{CommandStatus: status, Sequence: sequence},
				Tags:   Tags{0xFFFF: []byte(err.Error())},
			})
			continue
		}
		if callback, ok := s.pending[sequence]; ok {
			callback(packet)
		} else {
			s.receiveQueue <- packet
		}
	}
	_ = s.Close()
}

func (s *Session) Submit(ctx context.Context, packet Responsable) (resp Packet, err error) {
	sequence := s.NextSequence()
	WriteSequence(packet, sequence)
	if err = s.Send(packet); err != nil {
		return
	}
	returns := make(chan Packet, 1)
	s.pending[sequence] = func(resp Packet) { returns <- resp }
	defer delete(s.pending, sequence)
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case resp = <-returns:
	}
	return
}

func (s *Session) SubmitTimeout(packet Responsable, timeout time.Duration) (Packet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.Submit(ctx, packet)
}

func (s *Session) Send(packet Packet) (err error) {
	sequence := ReadSequence(packet)
	if sequence < 0 {
		err = ErrInvalidSequence
		return
	}
	if s.WriteTimeout > 0 {
		err = s.parent.SetWriteDeadline(time.Now().Add(s.WriteTimeout))
	}
	if err == nil {
		_, err = Marshal(s.parent, packet)
	}
	if err == io.EOF {
		err = ErrConnectionClosed
	}
	return
}

func (s *Session) EnquireLink(tick time.Duration, timeout time.Duration) {
	ticker := time.NewTicker(tick)
	var err error
	for range ticker.C {
		_, err = s.SubmitTimeout(new(EnquireLink), timeout)
		if err != nil {
			break
		}
	}
	ticker.Stop()
	close(s.receiveQueue)
}

func (s *Session) PDU() <-chan Packet {
	return s.receiveQueue
}

func (s *Session) Close() (err error) {
	_, err = s.SubmitTimeout(new(Unbind), time.Second)
	if err == nil {
		err = s.parent.Close()
	}
	if err == nil {
		close(s.receiveQueue)
	}
	return
}
