package pdu

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserDataHeader(t *testing.T) {
	mapping := map[string]UserDataHeader{
		"0500030C0201":   {0x00: []byte{12, 2, 1}},
		"060804F42E0201": {0x08: []byte{0xF4, 0x2E, 2, 1}},
	}
	concatenatedMapping := map[string]*ConcatenatedHeader{
		"0500030C0201":   {Reference: 12, TotalParts: 2, Sequence: 1},
		"060804F42E0201": {Reference: 62510, TotalParts: 2, Sequence: 1},
	}
	for packet, expected := range mapping {
		h := make(UserDataHeader)
		decoded, err := hex.DecodeString(packet)
		require.NoError(t, err)

		_, err = h.ReadFrom(bytes.NewReader(decoded))
		require.NoError(t, err)
		require.Equal(t, expected, h)
		require.Equal(t, len(decoded), h.Len())
		concatenated := h.ConcatenatedHeader()
		concatenated.Set(h)
		require.Equal(t, concatenatedMapping[packet], concatenated)
		require.Equal(t, len(decoded)-1, concatenated.Len())

		var buf bytes.Buffer
		_, err = h.WriteTo(&buf)
		require.NoError(t, err)
		require.Equal(t, decoded, buf.Bytes())
	}
	{
		h := make(UserDataHeader)
		_, err := h.ReadFrom(bytes.NewReader(nil))
		require.Error(t, err)
	}
	{
		h := UserDataHeader{0x08: []byte{0, 12, 2, 1}, 0x00: []byte{12, 2, 1}}
		var buf bytes.Buffer
		_, err := h.WriteTo(&buf)
		require.NoError(t, err)
	}
}

func TestUserDataHeader_ConcatenatedHeader(t *testing.T) {
	header := UserDataHeader{0x05: nil}
	require.Nil(t, header.ConcatenatedHeader())
}
