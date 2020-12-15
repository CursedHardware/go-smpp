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
	if err = readHeaderFrom(r, header); err != nil {
		return
	}
	if _, err = r.Read(make([]byte, header.CommandLength-16)); err != nil {
		err = ErrInvalidCommandLength
	}
	if t, ok := types[header.CommandID]; !ok {
		err = ErrInvalidCommandID
	} else {
		pdu = reflect.New(t).Interface()
		_, err = unmarshal(&buf, pdu)
	}
	return
}
