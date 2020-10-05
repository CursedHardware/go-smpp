package pdu

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserDataHeader(t *testing.T) {
	mapping := map[string]UserDataHeader{
		"0500030C0201":   {{ID: 0x00, Data: []byte{12, 2, 1}}},
		"060804F42E0201": {{ID: 0x08, Data: []byte{0xF4, 0x2E, 2, 1}}},
	}
	concatenatedMapping := map[string]*ConcatenatedHeader{
		"0500030C0201":   {Reference: 12, TotalParts: 2, Sequence: 1},
		"060804F42E0201": {Reference: 62510, TotalParts: 2, Sequence: 1},
	}
	for packet, expected := range mapping {
		h := UserDataHeader{}
		decoded, err := hex.DecodeString(packet)
		require.NoError(t, err)

		_, err = h.ReadFrom(bytes.NewReader(decoded))
		require.NoError(t, err)
		require.Equal(t, expected, h)
		require.Equal(t, len(decoded), h.Len())
		concatenated := h.ConcatenatedHeader()
		require.Equal(t, concatenatedMapping[packet], concatenated)
		require.Equal(t, len(decoded)-1, concatenated.Len())
		require.Equal(t, expected[0], concatenated.Element())

		var buf bytes.Buffer
		_, err = h.WriteTo(&buf)
		require.NoError(t, err)
		require.Equal(t, decoded, buf.Bytes())
	}
}

func TestInfoElement(t *testing.T) {
	element := InfoElement{ID: 0, Data: []byte{12, 2, 1}}
	require.Equal(t, "00:0C0201", element.String())

	var buf bytes.Buffer
	element.Data = make([]byte, 0x100)
	_, err := element.WriteTo(&buf)
	require.Error(t, err)
}

func TestUserDataHeader_ConcatenatedHeader(t *testing.T) {
	header := UserDataHeader{{ID: 0x05}}
	require.Nil(t, header.ConcatenatedHeader())
}
