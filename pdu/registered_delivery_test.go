package pdu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisteredDelivery(t *testing.T) {
	expected := byte(0b11110101)
	var delivery RegisteredDelivery
	_ = delivery.WriteByte(expected)
	require.Equal(t, delivery, RegisteredDelivery{
		MCDeliveryReceipt:           1,
		SMEOriginatedAcknowledgment: 1,
		IntermediateNotification:    true,
		Reserved:                    7,
	})
	c, _ := delivery.ReadByte()
	require.Equal(t, expected, c)
	require.Equal(t, "11110101", delivery.String())
}
