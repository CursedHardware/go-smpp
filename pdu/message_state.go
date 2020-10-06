package pdu

import (
	"strconv"
	"strings"
)

// MessageState see SMPP v5, section 4.7.15 (127p)
type MessageState byte

//goland:noinspection SpellCheckingInspection
var messageStateMap = []string{
	"scheduled",
	"enroute",
	"delivered",
	"expired",
	"deleted",
	"undeliverable",
	"accepted",
	"unknown",
	"rejected",
	"skipped",
}

func (m MessageState) String() string {
	if int(m) > len(messageStateMap) {
		return strconv.Itoa(int(m))
	}
	return strings.ToUpper(messageStateMap[m])
}
