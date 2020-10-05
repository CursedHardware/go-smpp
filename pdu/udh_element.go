package pdu

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type InfoElement struct {
	ID   byte
	Data []byte
}

func (e *InfoElement) ReadFrom(r io.Reader) (n int64, err error) {
	var length byte
	buf := bufio.NewReader(r)
	e.ID, err = buf.ReadByte()
	if err == nil {
		length, err = buf.ReadByte()
		e.Data = make([]byte, length)
	}
	if length > 0 {
		_, err = buf.Read(e.Data)
	}
	return
}

func (e InfoElement) WriteTo(w io.Writer) (n int64, err error) {
	if len(e.Data) > 0xFF {
		err = ErrDataTooLarge
		return
	}
	var buf bytes.Buffer
	buf.WriteByte(e.ID)
	buf.WriteByte(byte(len(e.Data)))
	buf.Write(e.Data)
	return buf.WriteTo(w)
}

func (e InfoElement) String() string {
	return fmt.Sprintf("%02X:%X", e.ID, e.Data)
}

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

func (h ConcatenatedHeader) Element() (element InfoElement) {
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.BigEndian, h)
	element.ID = 0x08
	element.Data = buf.Bytes()
	if h.Reference < 0xFF {
		element.ID = 0x00
		element.Data = element.Data[1:4]
	}
	return
}
