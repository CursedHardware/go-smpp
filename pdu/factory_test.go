package pdu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

//goland:noinspection SpellCheckingInspection
func TestCommandID(t *testing.T) {
	mapping := map[string]string{
		"SubmitSM":     "submit_sm",
		"SubmitSMResp": "submit_sm_resp",
	}
	for input, output := range mapping {
		require.Equal(t, output, toCommandIDName(input))
	}
	require.Equal(t, CommandID(0x00000004).String(), "submit_sm")
	require.Equal(t, CommandID(0xFFFFFFFF).String(), "FFFFFFFF")
}

//goland:noinspection SpellCheckingInspection
func TestCommandStatus(t *testing.T) {
	require.Equal(t, CommandStatus(0x00000000).String(), "ESME_ROK")
	require.Equal(t, CommandStatus(0xFFFFFFFF).String(), "FFFFFFFF")
}
