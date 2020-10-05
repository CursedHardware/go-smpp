package pdu

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Tags []Tag

type Tag struct {
	Tag  uint16
	Data []byte
}

func (t Tag) String() string {
	return fmt.Sprintf("%04X:%X", t.Tag, t.Data)
}

func (t Tags) FindTag(tag uint16) []byte {
	for _, item := range t {
		if item.Tag == tag {
			return item.Data
		}
	}
	return nil
}

func (t *Tags) ReadFrom(r io.Reader) (n int64, err error) {
	var fields Tags
	var values [2]uint16
	var data []byte
	for {
		err = binary.Read(r, binary.BigEndian, values[:])
		if err == nil {
			data = make([]byte, values[1])
			_, err = r.Read(data)
		}
		if err == nil {
			fields = append(fields, Tag{
				Tag:  values[0],
				Data: data,
			})
		}
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			break
		}
	}
	if len(fields) > 0 {
		*t = fields
	}
	return
}

func (t Tags) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer
	for _, field := range t {
		if length := len(field.Data); length == 0 {
			continue
		} else if length > 0xFFFF {
			err = ErrDataTooLarge
			return
		}
		_ = binary.Write(&buf, binary.BigEndian, field.Tag)
		_ = binary.Write(&buf, binary.BigEndian, uint16(len(field.Data)))
		buf.Write(field.Data)
	}
	return buf.WriteTo(w)
}
