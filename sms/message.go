package sms

import (
	"bufio"
	"io"
)

//goland:noinspection SpellCheckingInspection
func ParseMessage(r io.Reader, dir Direction) (smsc Address, packet interface{}, err error) {
	buf := bufio.NewReader(r)
	if smsc, err = ReadSMSCAddress(buf); err != nil {
		return
	}
	var t MessageType
	var failure bool
	peek, err := buf.Peek(2)
	if err != nil {
		return
	} else {
		t.Set(peek[0]&0b11, dir)
		failure = peek[1] > 0b001111111
	}
	switch t {
	case MessageTypeDeliver:
		packet = new(Deliver)
	case MessageTypeDeliverReport:
		packet = new(DeliverReport)
		if failure {
			packet = new(DeliverReportError)
		}
	case MessageTypeSubmit:
		packet = new(Submit)
	case MessageTypeSubmitReport:
		packet = new(SubmitReport)
		if failure {
			packet = new(SubmitReportError)
		}
	case MessageTypeStatusReport:
		packet = new(StatusReport)
	case MessageTypeCommand:
		packet = new(Command)
	default:
		return
	}
	_, err = unmarshal(buf, packet)
	return
}
