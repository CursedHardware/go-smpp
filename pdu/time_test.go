package pdu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {
	expectedList := map[string]string{
		"":                 "0001-01-01 00:00:00 +0000 UTC",
		"000101000000000+": "2000-01-01 00:00:00 +0000 +0000",
		"111019080000704-": "2011-10-19 08:00:00.7 -0100 -0100",
		"201020182347832+": "2020-10-20 18:23:47.8 +0800 +0800",
		"991231235959948+": "2099-12-31 23:59:59.9 +1200 +1200",
	}
	var timestamp Time
	for expected, formatted := range expectedList {
		require.NoError(t, timestamp.From(expected))
		require.Equal(t, formatted, timestamp.Time.String())
		require.Equal(t, expected, timestamp.String())
	}
	errorList := []string{
		"000101000000000",
	}
	for _, input := range errorList {
		require.Error(t, timestamp.From(input))
	}
}

func TestDuration(t *testing.T) {
	expectedList := map[string]string{
		"":                 "0s",
		"000007000000000R": "168h0m0s",
		"010203040506700R": "10276h5m6.7s",
		"991025033429000R": "875043h34m29s",
	}
	var duration Duration
	for expected, unix := range expectedList {
		require.NoError(t, duration.From(expected))
		require.Equal(t, unix, duration.Duration.String())
		require.Equal(t, expected, duration.String())
	}
	errorList := []string{
		"000101000000000",
	}
	for _, input := range errorList {
		require.Error(t, duration.From(input))
	}
}
