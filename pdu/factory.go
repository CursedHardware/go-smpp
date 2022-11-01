package pdu

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

var types = map[CommandID]reflect.Type{
	0x00000002: reflect.TypeOf(BindTransmitter{}),       // see SMPP v5, section 4.1.1.1 (56p)
	0x80000002: reflect.TypeOf(BindTransmitterResp{}),   // see SMPP v5, section 4.1.1.2 (57p)
	0x00000001: reflect.TypeOf(BindReceiver{}),          // see SMPP v5, section 4.1.1.3 (58p)
	0x80000001: reflect.TypeOf(BindReceiverResp{}),      // see SMPP v5, section 4.1.1.4 (59p)
	0x00000009: reflect.TypeOf(BindTransceiver{}),       // see SMPP v5, section 4.1.1.5 (59p)
	0x80000009: reflect.TypeOf(BindTransceiverResp{}),   // see SMPP v5, section 4.1.1.6 (60p)
	0x0000000B: reflect.TypeOf(Outbind{}),               // see SMPP v5, section 4.1.1.7 (61p)
	0x00000006: reflect.TypeOf(Unbind{}),                // see SMPP v5, section 4.1.1.8 (61p)
	0x80000006: reflect.TypeOf(UnbindResp{}),            // see SMPP v5, section 4.1.1.9 (62p)
	0x00000015: reflect.TypeOf(EnquireLink{}),           // see SMPP v5, section 4.1.2.1 (63p)
	0x80000015: reflect.TypeOf(EnquireLinkResp{}),       // see SMPP v5, section 4.1.2.2 (63p)
	0x00000102: reflect.TypeOf(AlertNotification{}),     // see SMPP v5, section 4.1.3.1 (64p)
	0x80000000: reflect.TypeOf(GenericNACK{}),           // see SMPP v5, section 4.1.4.1 (65p)
	0x00000004: reflect.TypeOf(SubmitSM{}),              // see SMPP v5, section 4.2.1.1 (66p)
	0x80000004: reflect.TypeOf(SubmitSMResp{}),          // see SMPP v5, section 4.2.1.2 (68p)
	0x00000103: reflect.TypeOf(DataSM{}),                // see SMPP v5, section 4.2.2.1 (69p)
	0x80000103: reflect.TypeOf(DataSMResp{}),            // see SMPP v5, section 4.2.2.2 (70p)
	0x00000021: reflect.TypeOf(SubmitMulti{}),           // see SMPP v5, section 4.2.3.1 (71p)
	0x80000021: reflect.TypeOf(SubmitMultiResp{}),       // see SMPP v5, section 4.2.3.2 (74p)
	0x00000005: reflect.TypeOf(DeliverSM{}),             // see SMPP v5, section 4.3.1.1 (85p)
	0x80000005: reflect.TypeOf(DeliverSMResp{}),         // see SMPP v5, section 4.3.1.1 (87p)
	0x00000112: reflect.TypeOf(BroadcastSM{}),           // see SMPP v5, section 4.4.1.1 (92p)
	0x80000112: reflect.TypeOf(BroadcastSMResp{}),       // see SMPP v5, section 4.4.1.2 (96p)
	0x00000008: reflect.TypeOf(CancelSM{}),              // see SMPP v5, section 4.5.1.1 (100p)
	0x80000008: reflect.TypeOf(CancelSMResp{}),          // see SMPP v5, section 4.5.1.2 (101p)
	0x00000003: reflect.TypeOf(QuerySM{}),               // see SMPP v5, section 4.5.2.1 (101p)
	0x80000003: reflect.TypeOf(QuerySMResp{}),           // see SMPP v5, section 4.5.2.2 (103p)
	0x00000007: reflect.TypeOf(ReplaceSM{}),             // see SMPP v5, section 4.5.3.1 (104p)
	0x80000007: reflect.TypeOf(ReplaceSMResp{}),         // see SMPP v5, section 4.5.3.2 (106p)
	0x00000111: reflect.TypeOf(QueryBroadcastSM{}),      // see SMPP v5, section 4.6.1.1 (107p)
	0x80000111: reflect.TypeOf(QueryBroadcastSMResp{}),  // see SMPP v5, section 4.6.1.3 (108p)
	0x00000113: reflect.TypeOf(CancelBroadcastSM{}),     // see SMPP v5, section 4.6.2.1 (110p)
	0x80000113: reflect.TypeOf(CancelBroadcastSMResp{}), // see SMPP v5, section 4.6.2.3 (112p)
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
	if t, ok := types[c]; ok {
		return toCommandIDName(t.Name())
	}
	return fmt.Sprintf("%08X", uint32(c))
}
