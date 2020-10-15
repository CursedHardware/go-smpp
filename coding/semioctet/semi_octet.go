package semioctet

import (
	"bytes"
	"io"
	"strconv"
)

func EncodeSemi(w io.Writer, chunks ...int) (n int64, err error) {
	digits := toDigits(chunks)
	var buf bytes.Buffer
	buf.Grow(len(digits) / 2)
	i, remain := 0, len(digits)
	for remain > 1 {
		buf.WriteByte(digits[i+1]<<4 | digits[i])
		i += 2
		remain -= 2
	}
	if remain > 0 {
		buf.WriteByte(0b11110000 | digits[i])
	}
	return buf.WriteTo(w)
}

func DecodeSemi(encoded []byte) (chunks []int) {
	var half byte
	for _, item := range encoded {
		half = item >> 4
		if half == 0b1111 {
			return append(chunks, int(item&0b1111))
		}
		chunks = append(chunks, int(item&0b1111*10+half))
	}
	return
}

func EncodeSemiAddress(w io.Writer, input string) (n int64, err error) {
	parsed, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return
	}
	return EncodeSemi(w, int(parsed))
}

func DecodeSemiAddress(encoded []byte) (output string) {
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

func toDigits(chunks []int) (digits []byte) {
	for _, chunk := range chunks {
		if chunk < 10 {
			digits = append(digits, 0)
		}
		for _, r := range strconv.Itoa(chunk) {
			digits = append(digits, byte(r-'0'))
		}
	}
	return
}
