package pdu

import (
	"bytes"
	"encoding/binary"
)

type ConcatenatedHeader struct {
	Reference  uint16
	TotalParts byte
	Sequence   byte
}

func (h ConcatenatedHeader) Len() int {
	if h.Reference < 0xFF {
		return 5
	}
	return 6
}

func (h ConcatenatedHeader) Set(udh UserDataHeader) {
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.BigEndian, h)
	id := byte(0x08)
	data := buf.Bytes()
	if h.Reference < 0xFF {
		id = 0x00
		data = data[1:4]
	}
	udh[id] = data
}
