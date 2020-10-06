package pdu

import (
	"bytes"
	"strings"
	"testing"

	"github.com/NiceLabs/go-smpp/coding"
	"github.com/stretchr/testify/require"
)

func TestShortMessage(t *testing.T) {
	var buf bytes.Buffer
	var message ShortMessage
	err := message.Compose("abc")
	require.NoError(t, err)
	_, err = message.Parse()
	require.NoError(t, err)
	err = message.Compose(strings.Repeat("abc", 54))
	require.Error(t, err)

	message.Message = make([]byte, 100)
	message.UDHeader = append(message.UDHeader, ConcatenatedHeader{
		Reference:  1,
		TotalParts: 1,
		Sequence:   1,
	}.Element())

	message.Message = make([]byte, MaxShortMessageLength+1)
	_, err = message.WriteTo(&buf)
	require.Error(t, err)

	message.DataCoding = coding.NoCoding
	parsed, err := message.Parse()
	require.NoError(t, err)
	require.NotEmpty(t, parsed)
}

func TestMessageState_String(t *testing.T) {
	require.Equal(t, "SCHEDULED", MessageState(0).String())
	require.Equal(t, "255", MessageState(0xFF).String())
}

//goland:noinspection SpellCheckingInspection
func TestComposeMultipartShortMessage(t *testing.T) {
	reference := uint16(0xFFFF)
	input := "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmno" +
		"pqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcde" +
		"fghijklmnopqrstuvwxyz1234123456789"
	expected := []ShortMessage{
		{
			Message:    []byte(input[:133]),
			UDHeader:   UserDataHeader{{ID: 0x08, Data: []byte{0xFF, 0xFF, 0x02, 0x01}}},
			DataCoding: coding.Latin1Coding,
		},
		{
			Message:    []byte(input[133:]),
			UDHeader:   UserDataHeader{{ID: 0x08, Data: []byte{0xFF, 0xFF, 0x02, 0x02}}},
			DataCoding: coding.Latin1Coding,
		},
	}
	messages, err := ComposeMultipartShortMessage(input, coding.Latin1Coding, reference)
	require.NoError(t, err)
	require.Equal(t, expected, messages)

	input = input[:10]
	expected = []ShortMessage{
		{Message: []byte(input), DataCoding: coding.Latin1Coding},
	}
	messages, err = ComposeMultipartShortMessage(input, coding.Latin1Coding, reference)
	require.NoError(t, err)
	require.Equal(t, expected, messages)
}

func TestComposeMultipartShortMessage_Error(t *testing.T) {
	input := make([]byte, 134*256)
	_, err := ComposeMultipartShortMessage(string(input), coding.NoCoding, 1)
	require.Error(t, err)
	_, err = ComposeMultipartShortMessage(string(input), coding.ASCIICoding, 1)
	require.Error(t, err)
	_, err = ComposeMultipartShortMessage(strings.Repeat("\xFF", 1000), coding.ASCIICoding, 1)
	require.Error(t, err)
}

func TestCombineMultipartDeliverSM(t *testing.T) {
	addDeliverSM := CombineMultipartDeliverSM(func([]*DeliverSM) {})
	addDeliverSM(&DeliverSM{
		Message: ShortMessage{Message: []byte(""), DataCoding: coding.Latin1Coding},
	})
	for i := 1; i < 3; i++ {
		header := UserDataHeader{{ID: 0x08, Data: []byte{0xFF, 0xFF, 0x02, byte(i)}}}
		addDeliverSM(&DeliverSM{
			Message: ShortMessage{Message: []byte(""), UDHeader: header, DataCoding: coding.Latin1Coding},
		})
	}
}
