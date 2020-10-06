package pdu

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
	"strconv"
)

func unmarshal(r io.Reader, packet interface{}) (n int64, err error) {
	buf := bufio.NewReader(r)
	v := reflect.ValueOf(packet)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		switch field := v.Field(i); field.Kind() {
		case reflect.String:
			var value string
			if value, err = readCString(buf); err == nil {
				field.SetString(value)
			}
		case reflect.Uint8:
			var value byte
			if value, err = buf.ReadByte(); err == nil {
				field.SetUint(uint64(value))
			}
		case reflect.Bool:
			var value byte
			if value, err = buf.ReadByte(); err == nil {
				field.SetBool(value == 1)
			}
		case reflect.Slice, reflect.Struct:
			switch v := (field.Addr().Interface()).(type) {
			case *Header:
				err = readHeaderFrom(buf, v)
			case io.ByteWriter:
				var value byte
				if value, err = buf.ReadByte(); err == nil {
					err = v.WriteByte(value)
				}
			case io.ReaderFrom:
				if m, ok := v.(*ShortMessage); ok {
					m.Prepare(packet)
				}
				_, err = v.ReadFrom(buf)
			}
		}
		n = int64(buf.Size())
		if err == io.EOF {
			err = nil
			return
		} else if err != nil {
			err = ErrUnmarshalPDUFailed
			return
		}
	}
	return
}

func Marshal(w io.Writer, packet interface{}) (n int64, err error) {
	var buf bytes.Buffer
	p := reflect.ValueOf(packet)
	if p.Kind() == reflect.Ptr {
		p = p.Elem()
	}
	for i := 0; i < p.NumField(); i++ {
		field := p.Field(i)
		switch field.Kind() {
		case reflect.String:
			writeCString(&buf, field.String())
		case reflect.Uint8:
			buf.WriteByte(byte(field.Uint()))
		case reflect.Bool:
			var value byte
			if field.Bool() {
				value = 1
			}
			buf.WriteByte(value)
		case reflect.Array, reflect.Slice, reflect.Struct:
			switch v := field.Addr().Interface().(type) {
			case *Header:
				var parsed uint64
				if value := p.Type().Field(i).Tag.Get(_ID); value != "" {
					parsed, err = strconv.ParseUint(value, 16, 32)
					v.CommandID = CommandID(parsed)
				}
				if err == nil && v.Sequence > 0 {
					_ = binary.Write(&buf, binary.BigEndian, v)
				} else {
					err = ErrInvalidSequence
				}
			case io.ByteReader:
				var value byte
				value, err = v.ReadByte()
				buf.WriteByte(value)
			case io.WriterTo:
				if m, ok := v.(*ShortMessage); ok {
					m.Prepare(packet)
				}
				_, err = v.WriteTo(&buf)
			}
		}
		if err != nil {
			return
		}
	}
	if p.Field(0).Type() == reflect.TypeOf(Header{}) {
		data := buf.Bytes()
		binary.BigEndian.PutUint32(data[0:4], uint32(buf.Len()))
	}
	return buf.WriteTo(w)
}
