package pdu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisteredDelivery(t *testing.T) {
	expected := byte(0b00010101)
	var delivery RegisteredDelivery
	_ = delivery.WriteByte(expected)
	c, _ := delivery.ReadByte()
	require.Equal(t, expected, c)
	require.Equal(t, "00010101", delivery.String())
}
