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
	t := p.Type()
	var validityPeriodFormat byte
	var parameterIndicator *ParameterIndicator
	for i := 0; i < p.NumField(); i++ {
		abbr := t.Field(i).Tag.Get("TP")
		if parameterIndicator != nil && !parameterIndicator.Has(abbr) {
			continue
		}
		field := p.Field(i).Addr().Interface()
		switch field := field.(type) {
		case *byte:
			*field, err = buf.ReadByte()
		case io.ByteWriter:
			var value byte
			if value, err = buf.ReadByte(); err == nil {
				err = field.WriteByte(value)
			}
			if setter, ok := field.(directionSetter); ok {
				switch t.Field(i).Tag.Get("DIR") {
				case "MT":
					setter.setDirection(MT)
				case "MO":
					setter.setDirection(MO)
				}
			}
		case *[]byte:
			var length byte
			if length, err = buf.ReadByte(); err == nil {
				*field = make([]byte, length)
				_, err = buf.Read(*field)
			}
		case io.ReaderFrom:
			_, err = field.ReadFrom(buf)
		case *ValidityPeriod:
			switch validityPeriodFormat {
			case 0b10:
				*field = new(Duration)
			case 0b11:
				*field = new(Time)
			}
			_, err = (*field).(io.ReaderFrom).ReadFrom(buf)
		}
		switch field := field.(type) {
		case *SubmitFlags:
			validityPeriodFormat = field.ValidityPeriodFormat
		case *ParameterIndicator:
			parameterIndicator = field
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

func unmarshalFlags(c byte, flags interface{}) (err error) {
	v := reflect.ValueOf(flags)
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
