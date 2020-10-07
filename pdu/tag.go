package pdu

import (
	"bytes"
	"encoding/binary"
	"io"
	"sort"
)

type Tags map[uint16][]byte

func (t *Tags) ReadFrom(r io.Reader) (n int64, err error) {
	var values [2]uint16
	var data []byte
	tags := make(Tags)
	for {
		err = binary.Read(r, binary.BigEndian, values[:])
		if err == nil {
			data = make([]byte, values[1])
			_, err = r.Read(data)
		}
		if err == nil {
			tags[values[0]] = data
		}
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			break
		}
	}
	if len(tags) > 0 {
		*t = tags
	}
	return
}

func (t Tags) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer
	var keys []uint16
	for tag := range t {
		keys = append(keys, tag)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	err = ErrDataTooLarge
	for _, tag := range keys {
		data := t[tag]
		length := len(data)
		if length == 0 {
			continue
		} else if length < 0xFFFF {
			_ = binary.Write(&buf, binary.BigEndian, tag)
			_ = binary.Write(&buf, binary.BigEndian, uint16(len(data)))
			buf.Write(data)
		} else {
			return
		}
	}
	return buf.WriteTo(w)
}
