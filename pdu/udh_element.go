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
	if data := buf.Bytes(); data[0] == 0 {
		udh[0x00] = data[1:4]
	} else {
		udh[0x08] = data
	}
}
