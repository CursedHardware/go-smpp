package pdu

import (
	"fmt"

	. "github.com/M2MGateway/go-smpp/coding"
)

func ComposeMultipartShortMessage(input string, coding DataCoding, reference uint16) (parts []ShortMessage, err error) {
	if coding.Splitter() == nil || coding.Encoding() == nil {
		err = ErrUnknownDataCoding
		return
	} else if coding.Splitter().Len(input) <= MaxShortMessageLength {
		var m ShortMessage
		m.DataCoding = coding
		m.Message, err = coding.Encoding().NewEncoder().Bytes([]byte(input))
		parts = []ShortMessage{m}
		return
	}
	header := ConcatenatedHeader{Reference: reference}
	segments := coding.Splitter().Split(input, MaxShortMessageLength-1-header.Len())
	if len(segments) > 0xFE {
		err = ErrMultipartTooMuch
		return
	}
	header.TotalParts = byte(len(segments))
	encoder := coding.Encoding().NewEncoder()
	part := ShortMessage{DataCoding: coding}
	for _, segment := range segments {
		encoder.Reset()
		part.UDHeader = make(UserDataHeader)
		if part.Message, err = encoder.Bytes([]byte(segment)); err != nil {
			return
		}
		header.Sequence++
		header.Set(part.UDHeader)
		parts = append(parts, part)
	}
	return
}

func CombineMultipartDeliverSM(on func([]*DeliverSM)) func(*DeliverSM) {
	registry := make(map[string][]*DeliverSM)
	isDone := func(id string, total byte) bool {
		for _, sm := range registry[id] {
			if sm != nil {
				total--
			}
		}
		return total == 0
	}
	return func(p *DeliverSM) {
		header := p.Message.UDHeader.ConcatenatedHeader()
		if header == nil {
			on([]*DeliverSM{p})
		} else {
			id := fmt.Sprint(
				p.SourceAddr.TON, p.SourceAddr.NPI, p.SourceAddr.No,
				p.DestAddr.TON, p.DestAddr.NPI, p.DestAddr.No,
				header.Reference,
			)
			if _, ok := registry[id]; !ok {
				registry[id] = make([]*DeliverSM, header.TotalParts)
			}
			registry[id][header.Sequence-1] = p
			if isDone(id, header.TotalParts) {
				on(registry[id])
				delete(registry, id)
			}
		}
	}
}
