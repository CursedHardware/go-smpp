package pdu

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

//goland:noinspection SpellCheckingInspection
func TestHeader(t *testing.T) {
	expectedList := map[string]Header{
		"00000010000000150000000000000007": {
			CommandLength: 16,
			CommandID:     0x00000015,
			Sequence:      7,
		},
		"000000100000FFFF0000FFFF0000FFFF": {
			CommandLength: 16,
			CommandID:     0x0000FFFF,
			CommandStatus: 0x0000FFFF,
			Sequence:      0x0000FFFF,
		},
	}
	var header Header
	for packet, expected := range expectedList {
		decoded, _ := hex.DecodeString(packet)
		err := readHeaderFrom(bytes.NewReader(decoded), &header)
		require.NoError(t, err)
		require.Equal(t, expected, header)
	}
	errorList := []string{
		"00000000000000000000000000000000",
		"00010001000000000000000000000000",
	}
	for _, packet := range errorList {
		decoded, _ := hex.DecodeString(packet)
		err := readHeaderFrom(bytes.NewReader(decoded), &header)
		require.Error(t, err)
	}
}

func TestSequence(t *testing.T) {
	ReadSequence(new(struct{}))
	ReadSequence(&DeliverSM{})
	_ = ReadCommandStatus(new(struct{}))
	_ = ReadCommandStatus(&DeliverSM{})
	WriteSequence(&DeliverSM{}, 0)
}
