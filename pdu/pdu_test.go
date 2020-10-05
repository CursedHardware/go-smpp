package pdu

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

//goland:noinspection SpellCheckingInspection
func TestReadPDU(t *testing.T) {
	failedList := []string{
		"00000000000000000000000000000000",
		"000000100000FFFF0000FFFF0000FFFF",
	}
	for _, packet := range failedList {
		decoded, err := hex.DecodeString(packet)
		require.NoError(t, err)
		_, err = ReadPDU(bytes.NewReader(decoded))
		require.Error(t, err)
	}
}
