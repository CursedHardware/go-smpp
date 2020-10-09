package pdu

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"sort"
)

type UserDataHeader map[byte][]byte

func (h UserDataHeader) Len() (length int) {
	if h == nil {
		return
	}
	length = 1
	for _, data := range h {
		length += 2
		length += len(data)
	}
	return
}

func (h *UserDataHeader) ReadFrom(r io.Reader) (n int64, err error) {
	buf := bufio.NewReader(r)
	header := make(UserDataHeader)
	length, err := buf.ReadByte()
	if err != nil {
		return
	}
	var id byte
	var data []byte
	for i := 0; i < int(length); {
		if id, err = buf.ReadByte(); err == nil {
			length, err = buf.ReadByte()
		}
		if length > 0 {
			data = make([]byte, length)
			_, err = buf.Read(data)
		}
		if err == nil {
			header[id] = data
		}
		i = buf.Size()
	}
	if len(header) > 0 {
		*h = header
	}
	return
}

func (h UserDataHeader) WriteTo(w io.Writer) (n int64, err error) {
	if h == nil {
		return
	}
	var keys []byte
	for id := range h {
		keys = append(keys, id)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	var buf bytes.Buffer
	buf.WriteByte(0)
	err = ErrDataTooLarge
	for _, id := range keys {
		data := h[id]
		if len(data) > 0xFF {
			return
		}
		buf.WriteByte(id)
		buf.WriteByte(byte(len(data)))
		buf.Write(data)
	}
	data := buf.Bytes()
	data[0] = byte(len(data)) - 1
	return buf.WriteTo(w)
}

func (h UserDataHeader) ConcatenatedHeader() *ConcatenatedHeader {
	for id, data := range h {
		switch id {
		case 0x00:
			return &ConcatenatedHeader{
				Reference:  uint16(data[0]),
				TotalParts: data[1],
				Sequence:   data[2],
			}
		case 0x08:
			return &ConcatenatedHeader{
				Reference:  binary.BigEndian.Uint16(data[0:2]),
				TotalParts: data[2],
				Sequence:   data[3],
			}
		}
	}
	return nil
}
