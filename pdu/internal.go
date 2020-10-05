package pdu

import (
	"bufio"
	"bytes"
)

func readCString(buf *bufio.Reader) (value string, err error) {
	value, err = buf.ReadString(0)
	if err == nil {
		value = value[0 : len(value)-1]
	}
	return
}

func writeCString(buf *bytes.Buffer, value string) {
	buf.WriteString(value)
	buf.WriteByte(0)
}
