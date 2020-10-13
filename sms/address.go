package sms

import (
	"bufio"
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
	if length == 0 {
		return
	}
	if kind, err = buf.ReadByte(); err != nil {
		return
	}
	p.NPI = kind & 0b1111
	p.TON = kind >> 4 & 0b111
	data := make([]byte, (length+1)/2)
	if _, err = buf.Read(data); err != nil {
		return
	}
	if p.TON != 0b101 {
		p.No = semioctet.DecodeSemiAddress(data)
	} else {
		data, err = gsm7bit.Packed.NewDecoder().Bytes(data)
		if err == nil {
			p.No = string(data)
		}
	}
	return
}

type SCAddress Address

func (p *SCAddress) ReadFrom(r io.Reader) (n int64, err error) {
	buf := bufio.NewReader(r)
	var length, kind byte
	if length, err = buf.ReadByte(); err != nil {
		return
	}
	if length == 0 {
		return
	}
	if kind, err = buf.ReadByte(); err != nil {
		return
	}
	p.NPI = kind & 0b1111
	p.TON = kind >> 4 & 0b111
	data := make([]byte, length-1)
	if _, err = buf.Read(data); err != nil {
		return
	}
	if p.TON != 0b101 {
		p.No = semioctet.DecodeSemiAddress(data)
	} else {
		data, err = gsm7bit.Packed.NewDecoder().Bytes(data)
		if err == nil {
			p.No = string(data)
		}
	}
	return
}
