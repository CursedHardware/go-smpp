package pdu

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Address struct {
	TON byte // see SMPP v5, section 4.7.1 (113p)
	NPI byte // see SMPP v5, section 4.7.2 (113p)
	No  string
}

func (p *Address) ReadFrom(r io.Reader) (n int64, err error) {
	buf := bufio.NewReader(r)
	p.TON, err = buf.ReadByte()
	if err == nil {
		p.NPI, err = buf.ReadByte()
	}
	if err == nil {
		p.No, err = readCString(buf)
	}
	return
}

func (p Address) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer
	buf.WriteByte(p.TON)
	buf.WriteByte(p.NPI)
	writeCString(&buf, p.No)
	return buf.WriteTo(w)
}

func (p Address) String() string {
	if p.TON == 1 && p.NPI == 1 && len(p.No) > 0 && p.No[0] != '+' {
		return "+" + p.No
	}
	return p.No
}

type DestinationAddresses struct {
	Addresses        []Address
	DistributionList []string
}

func (p *DestinationAddresses) ReadFrom(r io.Reader) (n int64, err error) {
	buf := bufio.NewReader(r)
	count, err := buf.ReadByte()
	if err != nil {
		err = ErrInvalidCommandLength
		return
	}
	*p = DestinationAddresses{}
	var destFlag byte
	var value string
	var address Address
	for i := byte(0); i < count; i++ {
		switch destFlag, _ = buf.ReadByte(); destFlag {
		case 1:
			if _, err = address.ReadFrom(buf); err == nil {
				p.Addresses = append(p.Addresses, address)
			}
		case 2:
			if value, err = readCString(buf); err == nil {
				p.DistributionList = append(p.DistributionList, value)
			}
		default:
			err = ErrInvalidDestFlag
			return
		}
		if err != nil {
			err = ErrInvalidCommandLength
			return
		}
	}
	return
}

func (p DestinationAddresses) WriteTo(w io.Writer) (n int64, err error) {
	length := len(p.Addresses) + len(p.DistributionList)
	if length > 0xFF {
		err = ErrInvalidDestCount
		return
	}
	var buf bytes.Buffer
	buf.WriteByte(byte(length))
	for _, address := range p.Addresses {
		buf.WriteByte(1)
		_, _ = address.WriteTo(&buf)
	}
	for _, distribution := range p.DistributionList {
		buf.WriteByte(2)
		writeCString(&buf, distribution)
	}
	return buf.WriteTo(w)
}

type UnsuccessfulRecords []UnsuccessfulRecord

type UnsuccessfulRecord struct {
	DestAddr        Address
	ErrorStatusCode CommandStatus
}

func (i UnsuccessfulRecord) String() string {
	return fmt.Sprintf("%s#%s", i.DestAddr, i.ErrorStatusCode)
}

func (p *UnsuccessfulRecords) ReadFrom(r io.Reader) (n int64, err error) {
	buf := bufio.NewReader(r)
	count, err := buf.ReadByte()
	if err != nil {
		err = ErrInvalidCommandLength
		return
	}
	items := UnsuccessfulRecords{}
	var item UnsuccessfulRecord
	for i := byte(0); i < count; i++ {
		_, err = item.DestAddr.ReadFrom(buf)
		if err == nil {
			err = binary.Read(buf, binary.BigEndian, &item.ErrorStatusCode)
		}
		if err != nil {
			err = ErrInvalidCommandLength
			return
		}
		items = append(items, item)
	}
	*p = items
	return
}

func (p UnsuccessfulRecords) WriteTo(w io.Writer) (n int64, err error) {
	if len(p) > 0xFF {
		err = ErrItemTooMany
		return
	}
	var buf bytes.Buffer
	buf.WriteByte(byte(len(p)))
	for _, item := range p {
		_, _ = item.DestAddr.WriteTo(&buf)
		_ = binary.Write(&buf, binary.BigEndian, item.ErrorStatusCode)
	}
	return buf.WriteTo(w)
}
