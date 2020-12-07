package pdu

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"io"
	"reflect"

	. "github.com/M2MGateway/go-smpp/coding"
)

type ShortMessage struct {
	DefaultMessageID byte // see SMPP v5, section 4.7.27 (134p)
	DataCoding       DataCoding
	UDHeader         UserDataHeader
	Message          []byte
}

func (p *ShortMessage) ReadFrom(r io.Reader) (n int64, err error) {
	buf := bufio.NewReader(r)
	if p.DataCoding != NoCoding {
		coding, _ := buf.ReadByte()
		p.DataCoding = DataCoding(coding)
	}
	p.DefaultMessageID, err = buf.ReadByte()
	if err == nil {
		var length byte
		if length, err = buf.ReadByte(); err == nil && p.UDHeader != nil {
			_, err = p.UDHeader.ReadFrom(buf)
		}
		if err == nil {
			p.Message = make([]byte, length-byte(p.UDHeader.Len()))
			_, err = buf.Read(p.Message)
		}
	}
	return
}

func (p ShortMessage) WriteTo(w io.Writer) (n int64, err error) {
	if len(p.Message) > MaxShortMessageLength {
		err = ErrShortMessageTooLarge
		return
	}
	var buf bytes.Buffer
	if p.DataCoding != NoCoding {
		buf.WriteByte(byte(p.DataCoding))
	}
	buf.WriteByte(p.DefaultMessageID)
	start := buf.Len()
	buf.WriteByte(0)
	_, err = p.UDHeader.WriteTo(&buf)
	if err != nil {
		return
	}
	buf.Write(p.Message)
	data := buf.Bytes()
	data[start] = byte(len(data) - 1 - start)
	return buf.WriteTo(w)
}

func (p *ShortMessage) Prepare(pdu interface{}) {
	if _, ok := pdu.(*ReplaceSM); ok {
		p.DataCoding = NoCoding
	} else if p.UDHeader == nil {
		v := reflect.ValueOf(pdu).Elem().FieldByName(_ESMClass)
		if v.IsValid() && v.Interface().(ESMClass).UDHIndicator {
			p.UDHeader = UserDataHeader{}
		}
	}
}

func (p *ShortMessage) Parse() (message string, err error) {
	encoder := p.DataCoding.Encoding()
	if encoder == nil {
		message = hex.EncodeToString(p.Message)
		return
	}
	decoded, err := encoder.NewDecoder().Bytes(p.Message)
	message = string(decoded)
	return
}

func (p *ShortMessage) Compose(input string) (err error) {
	coding := BestCoding(input)
	if coding.Splitter().Len(input) > MaxShortMessageLength {
		return ErrShortMessageTooLarge
	}
	message, err := coding.Encoding().NewEncoder().Bytes([]byte(input))
	if err == nil {
		p.DataCoding = coding
		p.Message = message
	}
	return
}
