package pdu

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

var types = map[CommandID]reflect.Type{}

var commandIDNames = map[CommandID]string{}

func init() {
	pduTypes := []interface{}{
		AlertNotification{}, GenericNACK{}, Outbind{},
		BindReceiver{}, BindReceiverResp{},
		BindTransceiver{}, BindTransceiverResp{},
		BindTransmitter{}, BindTransmitterResp{},
		BroadcastSM{}, BroadcastSMResp{},
		CancelBroadcastSM{}, CancelBroadcastSMResp{},
		CancelSM{}, CancelSMResp{},
		DataSM{}, DataSMResp{},
		DeliverSM{}, DeliverSMResp{},
		EnquireLink{}, EnquireLinkResp{},
		QueryBroadcastSM{}, QueryBroadcastSMResp{},
		QuerySM{}, QuerySMResp{},
		ReplaceSM{}, ReplaceSMResp{},
		SubmitMulti{}, SubmitMultiResp{},
		SubmitSM{}, SubmitSMResp{},
		Unbind{}, UnbindResp{},
	}
	var _parsed uint64
	var _id CommandID
	for _, pduType := range pduTypes {
		t := reflect.TypeOf(pduType)
		_parsed, _ = strconv.ParseUint(t.Field(0).Tag.Get(_ID), 16, 32)
		_id = CommandID(_parsed)
		types[_id] = t
		commandIDNames[_id] = toCommandIDName(t.Name())
	}
}

func toCommandIDName(name string) string {
	isUpper := unicode.IsUpper
	toLower := unicode.ToLower
	var b strings.Builder
	for i, r := range strings.ReplaceAll(name, "SM", "Sm") {
		if i > 0 && isUpper(r) {
			b.WriteRune('_')
		}
		b.WriteRune(toLower(r))
	}
	return b.String()
}

func (c CommandID) String() string {
	if name, ok := commandIDNames[c]; ok {
		return name
	}
	return fmt.Sprintf("%08X", uint32(c))
}
