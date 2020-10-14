package sms

import (
	"bufio"
	"bytes"
	"io"

	"github.com/NiceLabs/go-smpp/coding/gsm7bit"
	"github.com/NiceLabs/go-smpp/coding/semioctet"
)

type Address struct {
	NPI, TON byte
	No       string
}

func (p Address) MarshalBinary() (data []byte, err error) {
	var kind byte
	kind |= p.NPI & 0b1111
	kind |= p.TON & 0b111 << 4
	kind |= 1 << 7
	var buf bytes.Buffer
	buf.WriteByte(0x00)
	buf.WriteByte(kind)
	if p.TON != 0b101 {
		_, err = semioctet.EncodeSemiAddress(&buf, p.No)
	} else {
		_, err = gsm7bit.Packed.NewEncoder().Writer(&buf).Write([]byte(p.No))
	}
	data = buf.Bytes()
	data[0] = byte(len(data) - 2)
	return
}

func (p *Address) ReadFrom(r io.Reader) (n int64, err error) {
	buf := bufio.NewReader(r)
	var length, kind byte
	if length, err = buf.ReadByte(); err != nil || length == 0 {
		return
	}
	if kind, err = buf.ReadByte(); err != nil {
		return
	}
	p.NPI = kind & 0b1111
	p.TON = kind >> 4 & 0b111
	length = (length + 1) / 2
	data := make([]byte, length)
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

func (p *Address) WriteTo(w io.Writer) (n int64, err error) {
	if len(p.No) == 0 {
		_, err = w.Write([]byte{0})
		return
	}
	data, _ := p.MarshalBinary()
	data[0] *= 2
	if p.TON != 0b101 {
		data[0] -= 1
	}
	_, err = w.Write(data)
	return
}

type SCAddress Address

func (p *SCAddress) ReadFrom(r io.Reader) (n int64, err error) {
	buf := bufio.NewReader(r)
	var length, kind byte
	if length, err = buf.ReadByte(); err != nil || length == 0 {
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

func (p SCAddress) WriteTo(w io.Writer) (n int64, err error) {
	if len(p.No) == 0 {
		_, err = w.Write([]byte{0})
		return
	}
	data, _ := Address(p).MarshalBinary()
	data[0]++
	_, err = w.Write(data)
	return
}
