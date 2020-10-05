package pdu

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	var pdu interface{}

	pdu = new(DeliverSM)
	_, err := unmarshal(bytes.NewReader(nil), pdu)
	require.Error(t, err)

	pdu = &SubmitMulti{DestAddrList: DestinationAddresses{DistributionList: make([]string, 0x100)}}
	var buf bytes.Buffer
	_, err = Marshal(&buf, pdu)
	require.Error(t, err)

	pdu = &SubmitMulti{Header: Header{Sequence: -1}}
	_, err = Marshal(&buf, pdu)
	require.Error(t, err)
}
