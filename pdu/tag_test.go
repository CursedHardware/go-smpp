package pdu

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

//goland:noinspection SpellCheckingInspection
func TestTag(t *testing.T) {
	tlv := Tags{}
	require.Equal(t, tlv.FindTag(0x0007), []byte(nil))
	{
		_, err := tlv.ReadFrom(bytes.NewReader([]byte{0x00, 0x07, 0x00}))
		require.Error(t, err)
	}
	{
		_, err := tlv.ReadFrom(bytes.NewReader([]byte{0x00, 0x07, 0x00, 0x01, 0x5F}))
		require.NoError(t, err)
		require.Equal(t, tlv, Tags{{Tag: 0x0007, Data: []byte{0x5F}}})
	}
	{
		tlv[0].Data = []byte{}
		var buf bytes.Buffer
		_, err := tlv.WriteTo(&buf)
		require.NoError(t, err)
		require.Equal(t, 0, buf.Len())
		require.Equal(t, "[0007:]", fmt.Sprint(tlv))
	}
	{
		tlv[0].Data = make([]byte, 0x10000)
		var buf bytes.Buffer
		_, err := tlv.WriteTo(&buf)
		require.Error(t, err)
	}
	{
		tlv[0].Data = []byte{0x5F}
		require.Equal(t, "[0007:5F]", fmt.Sprint(tlv))
		require.Equal(t, tlv.FindTag(tlv[0].Tag), tlv[0].Data)
	}
	{
		tlv[0].Tag = 0xFFFF
		require.Equal(t, "[FFFF:5F]", fmt.Sprint(tlv))
	}
}
