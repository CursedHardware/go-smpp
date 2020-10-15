package semioctet

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSemiOctet(t *testing.T) {
	tests := map[string][]int{
		"1020304050": {1, 2, 3, 4, 5},
		"1122334455": {11, 22, 33, 44, 55},
	}
	for expected, input := range tests {
		decoded, err := hex.DecodeString(expected)
		require.NoError(t, err)
		var buf bytes.Buffer
		_, err = EncodeSemi(&buf, input...)
		require.NoError(t, err)
		require.Equal(t, decoded, buf.Bytes())
		require.Equal(t, input, DecodeSemi(decoded))
	}
	{
		require.Equal(t, []int{65, 53, 5}, DecodeSemi([]byte{0x56, 0x35, 0xF5}))
	}
}

func TestSemiOctetAddress(t *testing.T) {
	tests := map[string]string{
		"2143658709":   "1234567890",
		"2143658709F0": "12345678900",
	}
	for expected, input := range tests {
		decoded, err := hex.DecodeString(expected)
		require.NoError(t, err)
		var buf bytes.Buffer
		_, err = EncodeSemiAddress(&buf, input)
		require.NoError(t, err)
		require.Equal(t, decoded, buf.Bytes())
		require.Equal(t, input, DecodeSemiAddress(decoded))
	}
	{
		_, err := EncodeSemiAddress(nil, "ABC")
		require.Error(t, err)
	}
}
