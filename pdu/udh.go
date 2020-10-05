package pdu

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

// UserDataHeader: ID:Data
type UserDataHeader []InfoElement

func (h UserDataHeader) Len() (length int) {
	if h == nil {
		return 0
	}
	length = 1
	for _, element := range h {
		length += 2
		length += len(element.Data)
	}
	return
}

func (h *UserDataHeader) ReadFrom(r io.Reader) (n int64, err error) {
	buf := bufio.NewReader(r)
	header := UserDataHeader{}
	var element InfoElement
	length, err := buf.ReadByte()
	if err == nil {
		for i := 0; i < int(length); {
			if _, err = element.ReadFrom(buf); err == nil {
				header = append(header, element)
			}
			i = buf.Size()
		}
	}
	if len(header) > 0 {
		*h = header
	}
	return
}

func (h UserDataHeader) WriteTo(w io.Writer) (n int64, err error) {
	if h == nil {
		return 0, nil
	}
	var buf bytes.Buffer
	buf.WriteByte(0)
	for _, element := range h {
		_, _ = element.WriteTo(&buf)
	}
	data := buf.Bytes()
	data[0] = byte(len(data)) - 1
	return buf.WriteTo(w)
}

func (h UserDataHeader) ConcatenatedHeader() *ConcatenatedHeader {
	for _, element := range h {
		switch element.ID {
		case 0x00:
			return &ConcatenatedHeader{
				Reference:  uint16(element.Data[0]),
				TotalParts: element.Data[1],
				Sequence:   element.Data[2],
			}
		case 0x08:
			return &ConcatenatedHeader{
				Reference:  binary.BigEndian.Uint16(element.Data[0:2]),
				TotalParts: element.Data[2],
				Sequence:   element.Data[3],
			}
		}
	}
	return nil
}
