package pdu

import (
	"reflect"
)

func ReadSequence(packet any) int32 {
	if h := getHeader(packet); h != nil {
		return h.Sequence
	}
	return 0
}

func WriteSequence(packet any, sequence int32) {
	if h := getHeader(packet); h != nil {
		h.Sequence = sequence
	}
}

func ReadCommandStatus(packet any) CommandStatus {
	if h := getHeader(packet); h != nil {
		return h.CommandStatus
	}
	return 0
}

func getHeader(packet any) *Header {
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
