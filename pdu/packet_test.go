package pdu

import (
	"bytes"
	"encoding/hex"
	"testing"

	. "github.com/VoiceGateway/go-smpp/coding"
	"github.com/stretchr/testify/require"
)

//goland:noinspection SpellCheckingInspection
var (
	alice = Address{TON: 13, NPI: 15, No: "Alice"}
	bob   = Address{TON: 19, NPI: 7, No: "Bob"}
	empty = Address{TON: 23, NPI: 101, No: "empty"}
)

//goland:noinspection SpellCheckingInspection
var mapping = []struct {
	Packet         string
	Expected       interface{}
	Response       interface{}
	ResponsePacket string
}{
	{
		Packet: "0000003600000001000000000000000d73797374656d5f69645f66616b650070617373776f7264006f6e6c7900500d0f416c69636500",
		Expected: &BindReceiver{
			Header:       Header{54, 0x00000001, 0, 13},
			SystemID:     "system_id_fake",
			Password:     "password",
			SystemType:   "only",
			Version:      SMPPVersion50,
			AddressRange: alice,
		},
		Response: &BindReceiverResp{
			Header:   Header{31, 0x80000001, 0, 13},
			SystemID: "system_id_fake",
		},
		ResponsePacket: "0000001f80000001000000000000000d73797374656d5f69645f66616b6500",
	},
	{
		Packet: "00000024000000090000000000000001706f72742d31006d616e61676564000034000000",
		Expected: &BindTransceiver{
			Header:   Header{36, 0x00000009, 0, 1},
			SystemID: "port-1",
			Password: "managed",
			Version:  SMPPVersion34,
		},
		Response: &BindTransceiverResp{
			Header:   Header{23, 0x80000009, 0, 1},
			SystemID: "port-1",
		},
		ResponsePacket: "00000017800000090000000000000001706f72742d3100",
	},
	{
		Packet: "0000003600000002000000000000000d73797374656d5f69645f66616b650070617373776f7264006f6e6c7900501765656d70747900",
		Expected: &BindTransmitter{
			Header:       Header{54, 0x00000002, 0, 13},
			SystemID:     "system_id_fake",
			Password:     "password",
			SystemType:   "only",
			Version:      SMPPVersion50,
			AddressRange: empty,
		},
		Response: &BindTransmitterResp{
			Header:   Header{31, 0x80000002, 0, 13},
			SystemID: "system_id_fake",
		},
		ResponsePacket: "0000001f80000002000000000000000d73797374656d5f69645f66616b6500",
	},
	{
		Packet: "0000003200000112000000000000000d58585800010138363133383030313338303030006578616d706c6500000000010000",
		Expected: &BroadcastSM{
			Header:           Header{50, 0x00000112, 0, 13},
			ServiceType:      "XXX",
			SourceAddr:       Address{1, 1, "8613800138000"},
			MessageID:        "example",
			ReplaceIfPresent: true,
			DataCoding:       GSM7BitCoding,
		},
		Response: &BroadcastSMResp{
			Header:    Header{24, 0x80000112, 0, 13},
			MessageID: "example",
		},
		ResponsePacket: "0000001880000112000000000000000d6578616d706c6500",
	},
	{
		Packet: "0000002c00000113000000000000000d585858006578616d706c650001013836313338303031333830303000",
		Expected: &CancelBroadcastSM{
			Header:      Header{44, 0x00000113, 0, 13},
			MessageID:   "example",
			ServiceType: "XXX",
			SourceAddr:  Address{1, 1, "8613800138000"},
		},
		Response: &CancelBroadcastSMResp{
			Header: Header{16, 0x80000113, 0, 13},
		},
		ResponsePacket: "0000001080000113000000000000000d",
	},
	{
		Packet: "0000001e00000102000000000000000d0d0f416c696365001307426f6200",
		Expected: &AlertNotification{
			Header:     Header{30, 0x00000102, 0, 13},
			SourceAddr: alice,
			ESMEAddr:   bob,
		},
	},
	{
		Packet: "0000002300000008000000000000000d58585800000d0f416c696365001307426f6200",
		Expected: &CancelSM{
			Header:      Header{35, 0x00000008, 0, 13},
			ServiceType: "XXX",
			MessageID:   "",
			SourceAddr:  alice,
			DestAddr:    bob,
		},
		Response: &CancelSMResp{
			Header: Header{16, 0x80000008, 0, 13},
		},
		ResponsePacket: "0000001080000008000000000000000d",
	},
	{
		Packet: "0000001080000008000000000000000D",
		Expected: &CancelSMResp{
			Header: Header{16, 0x80000008, 0, 13},
		},
	},
	{
		Packet: "0000002a00000103000000000000000d616263000d0f416c696365001307426f62000d135b000700015f",
		Expected: &DataSM{
			Header:             Header{42, 0x00000103, 0, 13},
			ServiceType:        "abc",
			SourceAddr:         alice,
			DestAddr:           bob,
			ESMClass:           ESMClass{MessageType: 3, MessageMode: 1},
			RegisteredDelivery: RegisteredDelivery{MCDeliveryReceipt: 3, IntermediateNotification: true},
			DataCoding:         0b01011011,
			Tags:               Tags{0x0007: []byte{0x5F}},
		},
		Response:       &DataSMResp{Header: Header{17, 0x80000103, 0, 13}},
		ResponsePacket: "0000001180000103000000000000000d00",
	},
	{
		Packet: "000000BC00000005000000000000000958585800020131303031300002013000400000000000000800920500030503015C0A656C768475286237FF0C60A853EF4EE576F463A556DE590D63074EE48FDB884C4E1A52A167E58BE26216529E7406FF1A007F007F002000300030FF1A624B673A4E0A7F516D4191CF67E58BE2007F007F002000300031FF1A8D2662374F59989D007F007F002000300032FF1A5B9E65F68BDD8D39007F007F002000300033FF1A5E387528529E74064E1A",
		Expected: &DeliverSM{
			Header:      Header{188, 0x00000005, 0, 9},
			ServiceType: "XXX",
			SourceAddr:  Address{2, 1, "10010"},
			DestAddr:    Address{2, 1, "0"},
			ESMClass:    ESMClass{UDHIndicator: true},
			Message: ShortMessage{
				DataCoding: UCS2Coding,
				UDHeader:   UserDataHeader{0x00: []byte{0x05, 0x03, 0x01}},
				Message: []byte{
					0x5C, 0x0A, 0x65, 0x6C, 0x76, 0x84, 0x75, 0x28, 0x62, 0x37, 0xFF, 0x0C, 0x60, 0xA8, 0x53, 0xEF,
					0x4E, 0xE5, 0x76, 0xF4, 0x63, 0xA5, 0x56, 0xDE, 0x59, 0x0D, 0x63, 0x07, 0x4E, 0xE4, 0x8F, 0xDB,
					0x88, 0x4C, 0x4E, 0x1A, 0x52, 0xA1, 0x67, 0xE5, 0x8B, 0xE2, 0x62, 0x16, 0x52, 0x9E, 0x74, 0x06,
					0xFF, 0x1A, 0x00, 0x7F, 0x00, 0x7F, 0x00, 0x20, 0x00, 0x30, 0x00, 0x30, 0xFF, 0x1A, 0x62, 0x4B,
					0x67, 0x3A, 0x4E, 0x0A, 0x7F, 0x51, 0x6D, 0x41, 0x91, 0xCF, 0x67, 0xE5, 0x8B, 0xE2, 0x00, 0x7F,
					0x00, 0x7F, 0x00, 0x20, 0x00, 0x30, 0x00, 0x31, 0xFF, 0x1A, 0x8D, 0x26, 0x62, 0x37, 0x4F, 0x59,
					0x98, 0x9D, 0x00, 0x7F, 0x00, 0x7F, 0x00, 0x20, 0x00, 0x30, 0x00, 0x32, 0xFF, 0x1A, 0x5B, 0x9E,
					0x65, 0xF6, 0x8B, 0xDD, 0x8D, 0x39, 0x00, 0x7F, 0x00, 0x7F, 0x00, 0x20, 0x00, 0x30, 0x00, 0x33,
					0xFF, 0x1A, 0x5E, 0x38, 0x75, 0x28, 0x52, 0x9E, 0x74, 0x06, 0x4E, 0x1A,
				},
			},
		},
		Response:       &DeliverSMResp{Header: Header{17, 0x80000005, 0, 9}},
		ResponsePacket: "0000001180000005000000000000000900",
	},
	{
		Packet:         "00000010000000150000000000000007",
		Expected:       &EnquireLink{Header: Header{16, 0x00000015, 0, 7}},
		Response:       &EnquireLinkResp{Header: Header{16, 0x80000015, 0, 7}},
		ResponsePacket: "00000010800000150000000000000007",
	},
	{
		Packet:   "0000001080000000000000000000000D",
		Expected: &GenericNACK{Header: Header{16, 0x80000000, 0, 13}},
	},
	{
		Packet:   "000000240000000B000000000000000D696E76656E746F7279006970617373776F726400",
		Expected: &Outbind{Header: Header{36, 0x0000000B, 0, 13}, SystemID: "inventory", Password: "ipassword"},
	},
	{
		Packet: "0000002000000111000000000000000d6578616d706c65000d0f416c69636500",
		Expected: &QueryBroadcastSM{
			Header:     Header{32, 0x00000111, 0, 13},
			MessageID:  "example",
			SourceAddr: alice,
		},
		Response:       &QueryBroadcastSMResp{Header: Header{24, 0x80000111, 0, 13}, MessageID: "example"},
		ResponsePacket: "0000001880000111000000000000000d6578616d706c6500",
	},
	{
		Packet:         "0000001d00000003000000000000000d61776179000d0f416c69636500",
		Expected:       &QuerySM{Header: Header{29, 0x00000003, 0, 13}, MessageID: "away", SourceAddr: alice},
		Response:       &QuerySMResp{Header: Header{19, 0x80000003, 0, 13}},
		ResponsePacket: "0000001380000003000000000000000d000000",
	},
	{
		Packet: "0000002d00000007000000000000000d49445f486572000d0f416c6963650000001300096e6967687477697368",
		Expected: &ReplaceSM{
			Header:             Header{45, 0x00000007, 0, 13},
			MessageID:          "ID_Her",
			SourceAddr:         alice,
			RegisteredDelivery: RegisteredDelivery{MCDeliveryReceipt: 3, IntermediateNotification: true},
			Message:            ShortMessage{Message: []byte("nightwish"), DataCoding: NoCoding},
		},
		Response:       &ReplaceSMResp{Header: Header{16, 0x80000007, 0, 13}},
		ResponsePacket: "0000001080000007000000000000000D",
	},
	{
		Packet: "0000006d00000021000000000000000d585858000d0f416c6963650003010000426f623100024c6973743100024c69737432000d633d00001300080030006e006700681eaf0020006e00670068006900ea006e00670020006e0067006800691ec5006e00670020006e00671ea3",
		Expected: &SubmitMulti{
			Header:             Header{CommandLength: 109, CommandID: 0x00000021, Sequence: 13},
			ServiceType:        "XXX",
			SourceAddr:         alice,
			DestAddrList:       DestinationAddresses{Addresses: []Address{{No: "Bob1"}}, DistributionList: []string{"List1", "List2"}},
			ESMClass:           ESMClass{MessageType: 3, MessageMode: 1},
			ProtocolID:         99,
			PriorityFlag:       61,
			RegisteredDelivery: RegisteredDelivery{MCDeliveryReceipt: 3, IntermediateNotification: true},
			Message: ShortMessage{
				DataCoding: UCS2Coding,
				Message: []byte{
					0x00, 0x6E, 0x00, 0x67, 0x00, 0x68, 0x1E, 0xAF, 0x00, 0x20, 0x00, 0x6E, 0x00, 0x67, 0x00, 0x68,
					0x00, 0x69, 0x00, 0xEA, 0x00, 0x6E, 0x00, 0x67, 0x00, 0x20, 0x00, 0x6E, 0x00, 0x67, 0x00, 0x68,
					0x00, 0x69, 0x1E, 0xC5, 0x00, 0x6E, 0x00, 0x67, 0x00, 0x20, 0x00, 0x6E, 0x00, 0x67, 0x1E, 0xA3,
				},
			},
		},
		Response:       &SubmitMultiResp{Header: Header{18, 0x80000021, 0, 13}, UnsuccessfulSMEs: UnsuccessfulRecords{}},
		ResponsePacket: "0000001280000021000000000000000d0000",
	},
	{
		Packet: "0000003080000021000000000000000D666F6F7462616C6C00022621426F623100000000130000426F62320000000014",
		Expected: &SubmitMultiResp{
			Header:    Header{CommandLength: 48, CommandID: 0x80000021, Sequence: 13},
			MessageID: "football",
			UnsuccessfulSMEs: UnsuccessfulRecords{
				UnsuccessfulRecord{
					DestAddr:        Address{TON: 38, NPI: 33, No: "Bob1"},
					ErrorStatusCode: 19,
				},
				UnsuccessfulRecord{
					DestAddr:        Address{No: "Bob2"},
					ErrorStatusCode: 20,
				},
			},
		},
	},
	{
		Packet: "0000005c00000004000000000000000d585858000d0f416c696365001307426f62000d633d00001301080030006e006700681eaf0020006e00670068006900ea006e00670020006e0067006800691ec5006e00670020006e00671ea3",
		Expected: &SubmitSM{
			Header:             Header{CommandLength: 92, CommandID: 0x00000004, Sequence: 13},
			ServiceType:        "XXX",
			SourceAddr:         alice,
			DestAddr:           bob,
			ESMClass:           ESMClass{MessageType: 3, MessageMode: 1},
			ProtocolID:         99,
			PriorityFlag:       61,
			RegisteredDelivery: RegisteredDelivery{MCDeliveryReceipt: 3, IntermediateNotification: true},
			ReplaceIfPresent:   true,
			Message: ShortMessage{
				DataCoding: UCS2Coding,
				Message: []byte{
					0x00, 0x6E, 0x00, 0x67, 0x00, 0x68, 0x1E, 0xAF, 0x00, 0x20, 0x00, 0x6E, 0x00, 0x67, 0x00, 0x68,
					0x00, 0x69, 0x00, 0xEA, 0x00, 0x6E, 0x00, 0x67, 0x00, 0x20, 0x00, 0x6E, 0x00, 0x67, 0x00, 0x68,
					0x00, 0x69, 0x1E, 0xC5, 0x00, 0x6E, 0x00, 0x67, 0x00, 0x20, 0x00, 0x6E, 0x00, 0x67, 0x1E, 0xA3,
				},
			},
		},
		Response:       &SubmitSMResp{Header: Header{17, 0x80000004, 0, 13}},
		ResponsePacket: "0000001180000004000000000000000d00",
	},
	{
		Packet:         "0000001000000006000000000000000D",
		Expected:       &Unbind{Header: Header{16, 0x00000006, 0, 13}},
		Response:       &UnbindResp{Header: Header{16, 0x80000006, 0, 13}},
		ResponsePacket: "0000001080000006000000000000000D",
	},
}

func TestPacket(t *testing.T) {
	for _, sample := range mapping {
		decoded, err := hex.DecodeString(sample.Packet)
		require.NoError(t, err)

		var buf bytes.Buffer
		_, err = Marshal(&buf, sample.Expected)
		require.NoError(t, err, sample.Expected)
		require.Equal(t, decoded, buf.Bytes(), hex.EncodeToString(buf.Bytes()))

		parsed, err := ReadPDU(bytes.NewReader(decoded))
		require.NoError(t, err)
		require.NotNil(t, parsed)
		require.Equal(t, sample.Expected, parsed)

		if resp, ok := sample.Expected.(Responsable); ok {
			response := resp.Resp()
			require.NotNil(t, response)

			decoded, err = hex.DecodeString(sample.ResponsePacket)
			require.NoError(t, err)

			buf.Reset()
			_, err = Marshal(&buf, response)
			require.NoError(t, err)
			require.Equal(t, decoded, buf.Bytes(), hex.EncodeToString(buf.Bytes()))

			parsed, err = ReadPDU(&buf)
			require.NoError(t, err, resp)
			require.NotNil(t, parsed)
			require.Equal(t, sample.Response, parsed)
		}
	}
}
