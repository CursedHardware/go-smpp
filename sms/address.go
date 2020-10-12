package sms

import (
	"bufio"
	"bytes"
	"io"

	"github.com/NiceLabs/go-smpp/coding/gsm7bit"
	"github.com/NiceLabs/go-smpp/coding/semioctet"
)

type Address struct {
	NPI byte
	TON byte
	No  string
}

func (p *Address) ReadFrom(r io.Reader) (n int64, err error) {
	buf := bufio.NewReader(r)
	var length, kind byte
	if length, err = buf.ReadByte(); err != nil {
		return
	}
	if kind, err = buf.ReadByte(); err != nil {
		return
	}
	p.NPI = kind & 0b1111
	p.TON = kind >> 4 & 0b111
	length /= 2
	var data []byte
	if p.TON != 0b101 {
		data = make([]byte, length+1)
		if _, err = buf.Read(data); err != nil {
			return
		}
		p.No = semioctet.DecodeSemiAddress(data)
	} else {
		data = make([]byte, length)
		if _, err = buf.Read(data); err != nil {
			return
		}
		if data, err = gsm7bit.Packed.NewDecoder().Bytes(data); err == nil {
			p.No = string(data)
		}
	}
	return
}

func ReadSMSCAddress(r io.Reader) (address Address, err error) {
	var buf bytes.Buffer
	r = io.TeeReader(r, &buf)
	data := make([]byte, 1)
	if _, err = r.Read(data); err != nil {
		return
	}
	data = make([]byte, data[0])
	if _, err = r.Read(data); err != nil {
		return
	}
	_, err = address.ReadFrom(&buf)
	return
}
