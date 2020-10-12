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
	firstOctet, err := buf.ReadByte()
	if err != nil {
		return
	}
	_ = buf.UnreadByte()
	var t MessageType
	t.Set(firstOctet&0b11, dir)
	switch t {
	case MessageTypeDeliver:
		packet = new(Deliver)
	case MessageTypeDeliverReport:
		packet = new(DeliverReport)
	case MessageTypeSubmit:
		packet = new(Submit)
	case MessageTypeSubmitReport:
		packet = new(SubmitReport)
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
