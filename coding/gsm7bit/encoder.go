package gsm7bit

import (
	"bytes"

	"golang.org/x/text/transform"
)

type gsm7Encoder struct{}

func (e gsm7Encoder) Reset() { /* no needed */ }

func (e gsm7Encoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	if len(src) == 0 {
		return
	}
	septets, err := toSeptets(string(src))
	if err != nil {
		return
	}
	nDst = blocks(len(septets) * 7)
	if len(dst) < nDst {
		nDst = 0
		err = transform.ErrShortDst
		return
	}
	packSeptets(dst, septets)
	return
}

func packSeptets(dst []byte, septets []byte) {
	var index int
	var bit, item byte
	pack := func(c byte) {
		for i := 0; i < 7; i++ {
			dst[index] |= c >> i & 1 << bit
			bit++
			if bit == 8 {
				index++
				bit = 0
			}
		}
	}
	for _, c := range septets {
		item = c
		pack(c)
	}
	if 8-bit == 7 {
		pack(cr)
	} else if bit == 0 && item == cr {
		dst[index] = 0x00
		pack(cr)
	}
}

func toSeptets(input string) (septets []byte, err error) {
	var buf bytes.Buffer
	for _, r := range input {
		if v, ok := forwardLookup[r]; ok {
			buf.WriteByte(v)
		} else if v, ok := forwardEscapes[r]; ok {
			buf.WriteByte(esc)
			buf.WriteByte(v)
		} else {
			err = ErrInvalidCharacter
			return
		}
	}
	septets = buf.Bytes()
	return
}

func blocks(n int) (length int) {
	length = n / 8
	if n%8 != 0 {
		length += 1
	}
	return
}
