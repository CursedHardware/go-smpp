package sms

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
)

func unmarshal(buf *bufio.Reader, packet interface{}) (n int64, err error) {
	p := reflect.ValueOf(packet)
	if p.Kind() == reflect.Ptr {
		p = p.Elem()
	}
	var validityPeriodFormat byte
	for i := 0; i < p.NumField(); i++ {
		switch field := (p.Field(i).Addr().Interface()).(type) {
		case *byte:
			*field, err = buf.ReadByte()
		case *[]byte:
			length, err := buf.ReadByte()
			if err == nil {
				*field = make([]byte, length)
				_, err = buf.Read(*field)
			}
		case io.ReaderFrom:
			_, err = field.ReadFrom(buf)
		case *ValidityPeriod:
			switch validityPeriodFormat {
			case 0b10:
				var duration Duration
				if _, err = duration.ReadFrom(buf); err == nil {
					*field = duration
				}
			case 0b11:
				var time Time
				if _, err = time.ReadFrom(buf); err == nil {
					*field = time
				}
			}
		case *Flags, *SubmitFlags, *DeliverFlags:
			var value byte
			if value, err = buf.ReadByte(); err == nil {
				unmarshalFlags(value, field)
			}
			if flags, ok := field.(*SubmitFlags); ok {
				validityPeriodFormat = flags.ValidityPeriodFormat
			}
		}
		n = int64(buf.Size())
		if err != nil {
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
		case reflect.Uint8:
			buf.WriteByte(byte(field.Uint()))
		case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
			switch v := field.Addr().Interface().(type) {
			case *[]byte:
				buf.WriteByte(byte(len(*v)))
				buf.Write(*v)
			case io.ByteReader:
				var value byte
				value, err = v.ReadByte()
				buf.WriteByte(value)
			case io.WriterTo:
				_, err = v.WriteTo(&buf)
			}
		}
		if err != nil {
			return
		}
	}
	return buf.WriteTo(w)
}

func unmarshalFlags(c byte, packet interface{}) {
	v := reflect.ValueOf(packet)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i, bits := 0, byte(0); i < v.NumField(); i++ {
		switch field := (v.Field(i).Addr().Interface()).(type) {
		case *MessageType:
			field.Set(bits&0b11, field.Direction())
			bits += 2
		case *byte:
			*field = c >> bits & 0b11
			bits += 2
		case *bool:
			*field = c>>bits&0b1 == 1
			bits++
		}
	}
	return
}
