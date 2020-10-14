package sms

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

//goland:noinspection SpellCheckingInspection
func TestAddress(t *testing.T) {
	tests := map[string]Address{
		"00":                 {},
		"0B911604895626F9":   {NPI: 1, TON: 1, No: "61409865629"},
		"0ED1EDF27C1E3E97E7": {NPI: 1, TON: 5, No: "messages"},
		//"0DD1EDF27C1E3E9701": {NPI: 1, TON: 5, No: "message"},
		"0ED0D637396C7EBBCB": {NPI: 0, TON: 5, No: "Vodafone"},
	}
	for input, expected := range tests {
		var address Address
		decoded, err := hex.DecodeString(input)
		require.NoError(t, err)
		_, err = address.ReadFrom(bytes.NewReader(decoded))
		require.NoError(t, err)
		require.Equal(t, expected, address)
		var buf bytes.Buffer
		_, err = address.WriteTo(&buf)
		require.NoError(t, err)
		require.Equal(t, decoded, buf.Bytes())
	}
}

func TestSCAddress(t *testing.T) {
	tests := map[string]SCAddress{
		"00":                 {},
		"07911604895626F9":   {NPI: 1, TON: 1, No: "61409865629"},
		"08D1EDF27C1E3E97E7": {NPI: 1, TON: 5, No: "messages"},
		"08D0D637396C7EBBCB": {NPI: 0, TON: 5, No: "Vodafone"},
	}
	for input, expected := range tests {
		var address SCAddress
		decoded, err := hex.DecodeString(input)
		require.NoError(t, err)
		_, err = address.ReadFrom(bytes.NewReader(decoded))
		require.NoError(t, err)
		require.Equal(t, expected, address)
		var buf bytes.Buffer
		_, err = address.WriteTo(&buf)
		require.NoError(t, err)
		require.Equal(t, decoded, buf.Bytes())
	}
}

//goland:noinspection SpellCheckingInspection
func TestAddress_ErrorHandler(t *testing.T) {
	var address Address
	_, err := address.ReadFrom(bytes.NewReader([]byte{0x01}))
	require.Error(t, err)
	_, err = address.ReadFrom(bytes.NewReader([]byte{0x02, 0x02}))
	require.Error(t, err)
	var smsc SCAddress
	_, err = smsc.ReadFrom(bytes.NewReader([]byte{0x01}))
	require.Error(t, err)
	_, err = smsc.ReadFrom(bytes.NewReader([]byte{0x02, 0x02}))
	require.Error(t, err)
}
