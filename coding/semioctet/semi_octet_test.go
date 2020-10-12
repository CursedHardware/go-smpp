package semioctet

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSemiOctet(t *testing.T) {
	decoded, err := hex.DecodeString("9761989901F0")
	require.NoError(t, err)
	require.Equal(t, "79168999100", DecodeSemiAddress(decoded))
}
