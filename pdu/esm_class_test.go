package pdu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestESMClass(t *testing.T) {
	expected := byte(0b11001101)
	var esm ESMClass
	_ = esm.WriteByte(expected)
	c, _ := esm.ReadByte()
	require.Equal(t, esm, ESMClass{
		MessageMode:  1,
		MessageType:  3,
		UDHIndicator: true,
		ReplyPath:    true,
	})
	require.Equal(t, expected, c)
	require.Equal(t, "11001101", esm.String())
}
