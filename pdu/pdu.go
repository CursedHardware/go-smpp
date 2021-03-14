package pdu

import (
	"bytes"
	"io"
	"reflect"
)

func ReadPDU(r io.Reader) (pdu Packet, err error) {
	var buf bytes.Buffer
	r = io.TeeReader(r, &buf)
	header := new(Header)
	err = readHeaderFrom(r, header)
	if err != nil {
		return
	}
	if header.CommandLength > 16 {
		var n int
		n, err = r.Read(make([]byte, header.CommandLength-16))
		if err != nil {
			return
		} else if n != int(header.CommandLength-16) {
			err = ErrInvalidCommandLength
			return
		}
	}
	t, ok := types[header.CommandID]
	if !ok {
		err = ErrInvalidCommandID
		return
	}
	pdu = reflect.New(t).Interface().(Packet)
	_, err = unmarshal(&buf, pdu)
	return
}
