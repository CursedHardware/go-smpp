package pdu

import (
	"bytes"
	"io"
	"reflect"
)

func ReadPDU(r io.Reader) (pdu interface{}, err error) {
	var buf bytes.Buffer
	r = io.TeeReader(r, &buf)
	header := new(Header)
	err = readHeaderFrom(r, header)
	if err == nil && header.CommandLength > 16 {
		body := make([]byte, header.CommandLength-16)
		_, err = r.Read(body)
	}
	if err != nil {
		err = ErrInvalidCommandLength
		return
	}
	t, ok := types[header.CommandID]
	if !ok {
		err = ErrInvalidCommandID
		return
	}
	pdu = reflect.New(t).Interface()
	_, err = unmarshal(&buf, pdu)
	return
}
