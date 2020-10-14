package sms

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParameterIndicator(t *testing.T) {
	var indicator ParameterIndicator
	for _, abbr := range []string{"PID", "DCS", "UD", ""} {
		if !indicator.Has(abbr) {
			indicator.Set(abbr)
		}
	}
	value, err := indicator.ReadByte()
	require.NoError(t, err)
	require.Equal(t, byte(0b111), value)
	err = indicator.WriteByte(0b000)
	require.NoError(t, err)
	require.Equal(t, indicator, ParameterIndicator{
		ProtocolIdentifier: false,
		DataCoding:         false,
		UserData:           false,
	})
}

func TestFlags(t *testing.T) {
	tests := map[byte]interface{}{
		0b00000011: &Flags{
			MessageType: 0b110,
		},
		0b00111111: &DeliverFlags{
			MessageType:            0b110,
			MoreMessagesToSend:     true,
			ReplyPath:              true,
			UDHIndicator:           true,
			StatusReportIndication: true,
		},
		0b11111111: &SubmitFlags{
			MessageType:             0b110,
			RejectDuplicates:        true,
			ValidityPeriodFormat:    0b11,
			ReplyPath:               true,
			UserDataHeaderIndicator: true,
			StatusReportRequest:     true,
		},
	}
	for expected, flags := range tests {
		value, err := flags.(io.ByteReader).ReadByte()
		require.NoError(t, err)
		require.Equal(t, expected, value)
		err = flags.(io.ByteWriter).WriteByte(0)
		require.NoError(t, err)
		flags.(directionSetter).setDirection(MO)
	}
}

func TestFailureCause_Error(t *testing.T) {
	require.NotEmpty(t, FailureCause(0).Error())
	require.NotEmpty(t, FailureCause(0x80).Error())
}
