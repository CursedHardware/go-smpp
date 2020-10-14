package sms

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessageType_String(t *testing.T) {
	tests := []MessageType{
		MessageTypeDeliverReport,
		MessageTypeDeliver,
		MessageTypeSubmit,
		MessageTypeSubmitReport,
		MessageTypeCommand,
		MessageTypeStatusReport,
		MessageType(0xFF),
	}
	for _, kind := range tests {
		require.NotEmpty(t, kind.String())
	}
}
