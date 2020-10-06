package pdu

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	var pdu interface{}
	{
		pdu = new(DeliverSM)
		_, err := unmarshal(bytes.NewReader(nil), pdu)
		require.Error(t, err)
	}
	{
		pdu = new(SubmitSMResp)
		decoded, err := hex.DecodeString("00000010800000040000000b55104dc7")
		require.NoError(t, err)
		_, err = unmarshal(bytes.NewReader(decoded), pdu)
		require.NoError(t, err)
		var buf bytes.Buffer
		_, err = Marshal(&buf, pdu)
		require.NoError(t, err)
		require.Equal(t, decoded, buf.Bytes())
	}
	{
		pdu = &SubmitMulti{DestAddrList: DestinationAddresses{DistributionList: make([]string, 0x100)}}
		var buf bytes.Buffer
		_, err := Marshal(&buf, pdu)
		require.Error(t, err)
	}
	{
		var buf bytes.Buffer
		pdu = &SubmitMulti{Header: Header{Sequence: -1}}
		_, err := Marshal(&buf, pdu)
		require.Error(t, err)
	}
}
