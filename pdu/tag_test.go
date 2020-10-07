package pdu

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

//goland:noinspection SpellCheckingInspection
func TestTag(t *testing.T) {
	tlv := make(Tags)
	require.Equal(t, tlv[0x0007], []byte(nil))
	{
		_, err := tlv.ReadFrom(bytes.NewReader([]byte{0x00, 0x07, 0x00}))
		require.Error(t, err)
	}
	{
		expected := []byte{0x00, 0x07, 0x00, 0x01, 0x5F}
		_, err := tlv.ReadFrom(bytes.NewReader([]byte{0x00, 0x07, 0x00, 0x01, 0x5F}))
		require.NoError(t, err)
		require.Equal(t, tlv, Tags{0x0007: []byte{0x5F}})
		var buf bytes.Buffer
		_, err = tlv.WriteTo(&buf)
		require.NoError(t, err)
		require.Equal(t, expected, buf.Bytes())
	}
	{
		tlv[0x0007] = []byte{}
		var buf bytes.Buffer
		_, err := tlv.WriteTo(&buf)
		require.NoError(t, err)
		require.Equal(t, 0, buf.Len())
		tlv[0x0001] = make([]byte, 1)
		tlv[0x0002] = make([]byte, 1)
		tlv[0x0003] = make([]byte, 1)
		_, err = tlv.WriteTo(&buf)
		require.NoError(t, err)
		tlv[0x0007] = make([]byte, 0x10000)
		_, err = tlv.WriteTo(&buf)
		require.Error(t, err)
	}
}
