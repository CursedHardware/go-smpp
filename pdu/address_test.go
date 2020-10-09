package pdu

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddress_ReadFrom(t *testing.T) {
	mapping := map[string]Address{
		"315B7068616E746F6D537472696B6500": {TON: 0x31, NPI: 0x5B, No: "phantomStrike"},
		"5F0D7068616E746F6D4F7065726100":   {TON: 0x5F, NPI: 0x0D, No: "phantomOpera"},
	}
	for packet, expected := range mapping {
		var address Address
		decoded, err := hex.DecodeString(packet)
		require.NoError(t, err)

		_, err = address.ReadFrom(bytes.NewReader(decoded))
		require.NoError(t, err)
		require.Equal(t, expected, address)
		require.Equal(t, expected.No, address.String())

		var buf bytes.Buffer
		_, err = address.WriteTo(&buf)
		require.NoError(t, err)
		require.Equal(t, decoded, buf.Bytes())
	}
}

func TestAddress_String(t *testing.T) {
	mapping := map[string]Address{
		"+15417543010": {TON: 1, NPI: 1, No: "15417543010"},
	}
	for expected, address := range mapping {
		require.Equal(t, expected, address.String())
	}
}

func TestDestinationAddresses(t *testing.T) {
	mapping := map[string]DestinationAddresses{
		"03010000426F623100024C6973743100024C6973743200": {
			Addresses:        []Address{{No: "Bob1"}},
			DistributionList: []string{"List1", "List2"},
		},
	}
	for packet, expected := range mapping {
		addresses := DestinationAddresses{}
		decoded, err := hex.DecodeString(packet)
		require.NoError(t, err)

		_, err = addresses.ReadFrom(bytes.NewReader(decoded))
		require.NoError(t, err)
		require.Equal(t, expected, addresses)

		var buf bytes.Buffer
		_, err = addresses.WriteTo(&buf)
		require.NoError(t, err)
		require.Equal(t, decoded, buf.Bytes())
	}
}

func TestDestinationAddresses_ReadFrom(t *testing.T) {
	var addresses DestinationAddresses
	_, err := addresses.ReadFrom(bytes.NewReader([]byte{0xFF, 0x02}))
	require.Error(t, err)
	_, err = addresses.ReadFrom(bytes.NewReader([]byte{0xFF, 0x03}))
	require.Error(t, err)
	_, err = addresses.ReadFrom(bytes.NewReader(nil))
	require.Error(t, err)
	addresses.DistributionList = make([]string, 0x100)
	var buf bytes.Buffer
	_, err = addresses.WriteTo(&buf)
	require.Error(t, err)
}

func TestUnsuccessfulRecords(t *testing.T) {
	packet := "022621426F623100000000130000426F62320000000014"
	parsed := UnsuccessfulRecords{
		UnsuccessfulRecord{ErrorStatusCode: 19, DestAddr: Address{TON: 38, NPI: 33, No: "Bob1"}},
		UnsuccessfulRecord{ErrorStatusCode: 20, DestAddr: Address{No: "Bob2"}},
	}
	stringify := "[Bob1#ESME_RREPLACEFAIL Bob2#ESME_RMSGQFUL]"
	smes := UnsuccessfulRecords{}
	decoded, err := hex.DecodeString(packet)
	require.NoError(t, err)

	_, err = smes.ReadFrom(bytes.NewReader(decoded))
	require.NoError(t, err)
	require.Equal(t, parsed, smes)
	require.Equal(t, stringify, fmt.Sprint(smes))

	var buf bytes.Buffer
	_, err = smes.WriteTo(&buf)
	require.NoError(t, err)
	require.Equal(t, decoded, buf.Bytes())

	smes = make([]UnsuccessfulRecord, 0x100)
	_, err = smes.WriteTo(&buf)
	require.Error(t, err)
}

func TestUnsuccessfulRecords_ReadFrom(t *testing.T) {
	sme := UnsuccessfulRecords{}
	_, err := sme.ReadFrom(bytes.NewReader([]byte{0xFF}))
	require.Error(t, err)
	_, err = sme.ReadFrom(bytes.NewReader(nil))
	require.Error(t, err)
}
