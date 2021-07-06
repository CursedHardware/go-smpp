package main

import (
	"fmt"

	"github.com/M2MGateway/go-smpp/coding"
	"github.com/M2MGateway/go-smpp/pdu"
)

func getSegments(text string) string {
	var length int
	if encoding := coding.GSM7BitCoding; encoding.Validate(text) {
		length = encoding.Splitter().Len(text)
	} else {
		length = coding.UCS2Coding.Splitter().Len(text)
	}
	segments := float64(1)
	if length > pdu.MaxShortMessageLength {
		segments = float64(length) / float64(pdu.MaxShortMessageLength-7)
	}
	return fmt.Sprintf("%.2f", segments)
}
