package sms

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"reflect"
)

//goland:noinspection SpellCheckingInspection
func Unmarshal(r io.Reader) (packet any, err error) {
	buf := bufio.NewReader(r)
	kind, failure, err := getType(buf)
	if err != nil {
		return
	}
	switch {
	case kind == MessageTypeDeliver:
		packet = new(Deliver)
	case kind == MessageTypeDeliverReport && failure:
		packet = new(DeliverReportError)
	case kind == MessageTypeDeliverReport:
		packet = new(DeliverReport)
	case kind == MessageTypeSubmit:
		packet = new(Submit)
	case kind == MessageTypeSubmitReport && failure:
		packet = new(SubmitReportError)
	case kind == MessageTypeSubmitReport:
		packet = new(SubmitReport)
	case kind == MessageTypeStatusReport:
		packet = new(StatusReport)
	case kind == MessageTypeCommand:
		packet = new(Command)
	default:
		err = errors.New(kind.String())
		return
	}
	_, err = unmarshal(buf, packet)
	return
}

func unmarshal(buf *bufio.Reader, packet any) (n int64, err error) {
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
		case *any:
			switch {
			case abbr == "VP" && validityPeriodFormat == 0b01:
				var duration EnhancedDuration
				_, err = duration.ReadFrom(buf)
				*field = duration
			case abbr == "VP" && validityPeriodFormat == 0b10:
				var duration Duration
				_, err = duration.ReadFrom(buf)
				*field = duration
			case abbr == "VP" && validityPeriodFormat == 0b11:
				var time Time
				_, err = time.ReadFrom(buf)
				*field = time
			}
		}
		switch field := field.(type) {
		case *SubmitFlags:
			validityPeriodFormat = field.ValidityPeriodFormat
		case *ParameterIndicator:
			parameterIndicator = field
		}
		if err != nil {
			return
		}
	}
	return
}

func Marshal(w io.Writer, packet any) (n int64, err error) {
	p := reflect.ValueOf(packet)
	if p.Kind() == reflect.Ptr {
		p = p.Elem()
	}
	t := p.Type()
	var validityPeriodFormat byte
	var parameterIndicator ParameterIndicator
	for i := 0; i < p.NumField(); i++ {
		parameterIndicator.Set(t.Field(i).Tag.Get("TP"))
		switch field := p.Field(i).Addr().Interface().(type) {
		case *any:
			switch (*field).(type) {
			case EnhancedDuration:
				validityPeriodFormat = 0b01
			case Duration:
				validityPeriodFormat = 0b10
			case Time:
				validityPeriodFormat = 0b11
			}
		}
	}
	var buf bytes.Buffer
	for i := 0; i < p.NumField(); i++ {
		switch field := p.Field(i).Addr().Interface().(type) {
		case *byte:
			buf.WriteByte(*field)
		case *[]byte:
			length := len(*field)
			buf.WriteByte(byte(length))
			buf.Write(bytes.TrimRight(*field, "\x00"))
		case io.ByteReader:
			if flags, ok := field.(*SubmitFlags); ok {
				flags.ValidityPeriodFormat = validityPeriodFormat
			}
			var value byte
			if value, err = field.ReadByte(); err == nil {
				buf.WriteByte(value)
			}
		case io.WriterTo:
			_, err = field.WriteTo(&buf)
		case *any:
			switch field := (*field).(type) {
			case EnhancedDuration:
				_, err = field.WriteTo(&buf)
			case Duration:
				_, err = field.WriteTo(&buf)
			case Time:
				_, err = field.WriteTo(&buf)
			}
		}
		if err != nil {
			return
		}
	}
	return buf.WriteTo(w)
}

func getType(buf *bufio.Reader) (kind MessageType, failure bool, err error) {
	var peek []byte
	if peek, err = buf.Peek(1); err != nil {
		return
	}
	length := int(peek[0])
	if peek, err = buf.Peek(length + 3); err != nil {
		return
	}
	var dir Direction
	if length == 0 {
		dir = MO
	}
	kind.Set(peek[length+1]&0b11, dir)
	failure = peek[length+2] > 0b001111111
	return
}

func unmarshalFlags(c byte, flags any) (err error) {
	v := reflect.ValueOf(flags)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	var b byte
	for i, bits := 0, byte(0); i < v.NumField(); i++ {
		b = c >> bits
		switch field := v.Field(i).Addr().Interface().(type) {
		case *MessageType:
			field.Set(b&0b11, field.Direction())
			bits += 2
		case *byte:
			*field = b & 0b11
			bits += 2
		case *bool:
			*field = b&0b1 == 1
			bits++
		}
	}
	return
}

func marshalFlags(flags any) (c byte, err error) {
	v := reflect.ValueOf(flags)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i, bits := 0, byte(0); i < v.NumField(); i++ {
		switch field := (v.Field(i).Interface()).(type) {
		case MessageType:
			c |= field.Type() << bits
			bits += 2
		case byte:
			c |= (field & 0b11) << bits
			bits += 2
		case bool:
			if field {
				c |= 1 << bits
			}
			bits++
		}
	}
	return
}
