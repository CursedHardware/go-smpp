package pdu

import (
	"reflect"
)

func ReadSequence(packet interface{}) int32 {
	if h := getHeader(packet); h != nil {
		return h.Sequence
	}
	return 0
}

func WriteSequence(packet interface{}, sequence int32) {
	if h := getHeader(packet); h != nil {
		h.Sequence = sequence
	}
}

func ReadCommandStatus(packet interface{}) CommandStatus {
	if h := getHeader(packet); h != nil {
		return h.CommandStatus
	}
	return 0
}

func getHeader(packet interface{}) *Header {
	p := reflect.ValueOf(packet)
	if p.Kind() == reflect.Ptr {
		p = p.Elem()
	}
	for i := 0; i < p.NumField(); i++ {
		field := p.Field(i)
		if h, ok := field.Addr().Interface().(*Header); ok {
			return h
		}
	}
	return nil
}
