package sms

import (
	"bufio"
	"errors"
	"io"
)

//goland:noinspection SpellCheckingInspection
func ParseMessage(r io.Reader, dir Direction) (packet interface{}, err error) {
	buf := bufio.NewReader(r)
	kind, failure, err := getType(buf, dir)
	if err != nil {
		return
	}
	switch {
	case kind == MessageTypeDeliver:
		packet = new(Deliver)
	case kind == MessageTypeDeliverReport && failure:
		packet = new(DeliverReportError)
	case kind == MessageTypeDeliverReport:
		packet = new(DeliverReport)
	case kind == MessageTypeSubmit:
		packet = new(Submit)
	case kind == MessageTypeSubmitReport && failure:
		packet = new(SubmitReportError)
	case kind == MessageTypeSubmitReport:
		packet = new(SubmitReport)
	case kind == MessageTypeStatusReport:
		packet = new(StatusReport)
	case kind == MessageTypeCommand:
		packet = new(Command)
	default:
		err = errors.New(kind.String())
		return
	}
	_, err = unmarshal(buf, packet)
	return
}

func getType(buf *bufio.Reader, dir Direction) (kind MessageType, failure bool, err error) {
	peek, err := buf.Peek(1)
	if err != nil {
		return
	}
	start := int(peek[0])
	peek, err = buf.Peek(start + 3)
	if err != nil {
		return
	} else {
		kind.Set(peek[start+1]&0b11, dir)
		failure = peek[start+2] > 0b001111111
	}
	return
}
