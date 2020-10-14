package semioctet

import (
	"bytes"
	"io"
	"strconv"
)

func EncodeSemi(w io.Writer, chunks ...int) (n int64, err error) {
	digits := make([]byte, 0, len(chunks))
	for _, c := range chunks {
		var bucket []byte
		if c < 10 {
			digits = append(digits, 0)
		}
		for c > 0 {
			d := c % 10
			bucket = append(bucket, byte(d))
			c = (c - d) / 10
		}
		for i := range bucket {
			digits = append(digits, bucket[len(bucket)-1-i])
		}
	}
	var buf bytes.Buffer
	buf.Grow(len(digits)/2 + 1)
	for i := 0; i < len(digits); i += 2 {
		if len(digits)-i < 2 {
			buf.WriteByte(0b11110000 | digits[i])
			return buf.WriteTo(w)
		}
		buf.WriteByte(digits[i+1]<<4 | digits[i])
	}
	return buf.WriteTo(w)
}

func DecodeSemi(encoded []byte) []byte {
	chunks := make([]byte, 0, len(encoded)*2)
	var half byte
	for _, item := range encoded {
		half = item >> 4
		if half == 0b1111 {
			return append(chunks, item&0b1111)
		}
		chunks = append(chunks, item&0b1111*10+half)
	}
	return chunks
}

func EncodeSemiAddress(w io.Writer, content string) (n int64, err error) {
	phone, err := strconv.ParseUint(content, 10, 64)
	if err != nil {
		return
	}
	return EncodeSemi(w, int(phone))
}

func DecodeSemiAddress(encoded []byte) string {
	var buf bytes.Buffer
	var half byte
	for _, item := range encoded {
		half = item & 0b1111
		buf.WriteByte('0' + half)
		if half = item >> 4; half != 0b1111 {
			buf.WriteByte('0' + half)
		}
	}
	return buf.String()
}
