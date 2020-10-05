package pdu

import "fmt"

// ESMClass see SMPP v5, section 4.7.12 (125p)
type ESMClass struct {
	MessageMode, MessageType   byte
	SetReplyPath, UDHIndicator bool
}

func (e ESMClass) ReadByte() (c byte, err error) {
	c |= e.MessageMode & 0b11
	c |= e.MessageType & 0b1111 << 2
	if e.UDHIndicator {
		c |= 1 << 6
	}
	if e.SetReplyPath {
		c |= 1 << 7
	}
	return
}

func (e *ESMClass) WriteByte(c byte) error {
	e.MessageMode = c & 0b11
	e.MessageType = c >> 2 & 0b1111
	e.UDHIndicator = c>>6&0b1 == 1
	e.SetReplyPath = c>>7&0b1 == 1
	return nil
}

func (e ESMClass) String() string {
	c, _ := e.ReadByte()
	return fmt.Sprintf("%08b", c)
}
