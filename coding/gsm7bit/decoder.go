package gsm7bit

import (
	"bytes"

	"golang.org/x/text/transform"
)

type gsm7Decoder struct{}

func (d gsm7Decoder) Reset() { /* no needed */ }

func (d gsm7Decoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	if len(src) == 0 {
		return
	}
	var buf bytes.Buffer
	septets := unpackSeptets(src)
	err = ErrInvalidByte
	for i, septet := 0, byte(0); i < len(septets); i++ {
		septet = septets[i]
		if septet <= 0x7F && septet != esc {
			buf.WriteRune(reverseLookup[septet])
		} else {
			i++
			if i >= len(septets) {
				return
			}
			r, ok := reverseEscapes[septets[i]]
			if !ok {
				return
			}
			buf.WriteRune(r)
		}
	}
	err = nil
	nDst = buf.Len()
	if len(dst) < nDst {
		nDst = 0
		err = transform.ErrShortDst
	} else {
		decoded := buf.Bytes()
		if n := len(decoded); n > 2 && (decoded[n-1] == cr || decoded[n-2] == cr) {
			nDst--
		}
		copy(dst, decoded)
	}
	return
}

func unpackSeptets(septets []byte) []byte {
	var septet, bit byte = 0, 0
	var buf bytes.Buffer
	buf.Grow(len(septets))
	for _, octet := range septets {
		for i := 0; i < 8; i++ {
			septet |= octet >> i & 1 << bit
			bit++
			if bit == 7 {
				buf.WriteByte(septet)
				septet = 0
				bit = 0
			}
		}
	}
	return buf.Bytes()
}
