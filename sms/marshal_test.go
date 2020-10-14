package sms

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getType(t *testing.T) {
	_, _, err := getType(bufio.NewReader(bytes.NewReader(nil)))
	require.Error(t, err)
	_, _, err = getType(bufio.NewReader(bytes.NewReader([]byte{0x01})))
	require.Error(t, err)
}
