package sms

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {
	expected := "2017-08-31T11:21:54+08:00"
	var timestamp Time
	decoded, err := hex.DecodeString("71801311124523")
	require.NoError(t, err)
	_, err = timestamp.ReadFrom(bytes.NewReader(decoded))
	require.NoError(t, err)
	require.Equal(t, expected, timestamp.Format(time.RFC3339))
	var buf bytes.Buffer
	_, err = timestamp.WriteTo(&buf)
	require.NoError(t, err)
	require.Equal(t, decoded, buf.Bytes())
}

func TestDuration(t *testing.T) {
	tests := map[string]Duration{
		"00": {Duration: 5 * time.Minute},
		"83": {Duration: 11 * time.Hour},
		"A5": {Duration: 23 * time.Hour},
		"C3": {Duration: 29 * 24 * time.Hour},
		"FE": {Duration: 62 * 7 * 24 * time.Hour},
		"FF": {Duration: 63 * 7 * 24 * time.Hour},
	}
	for input, expected := range tests {
		var duration Duration
		decoded, err := hex.DecodeString(input)
		require.NoError(t, err)
		_, err = duration.ReadFrom(bytes.NewReader(decoded))
		require.NoError(t, err)
		require.Equal(t, expected, duration)
		var buf bytes.Buffer
		_, err = duration.WriteTo(&buf)
		require.NoError(t, err)
		require.Equal(t, decoded, buf.Bytes())
	}
}

func TestEnhancedDuration(t *testing.T) {
	tests := map[string]EnhancedDuration{
		"00000000000000": {},
		"01000000000000": {Indicator: 0b001, Duration: 5 * time.Minute},
		"01010000000000": {Indicator: 0b001, Duration: 10 * time.Minute},
		"01FE0000000000": {Indicator: 0b001, Duration: 62 * 7 * 24 * time.Hour},
		"01FF0000000000": {Indicator: 0b001, Duration: 63 * 7 * 24 * time.Hour},
		"02FF0000000000": {Indicator: 0b010, Duration: 255 * time.Second},
		"03302154000000": {Indicator: 0b011, Duration: 3*time.Hour + 12*time.Minute + 45*time.Second},
	}
	for input, expected := range tests {
		var duration EnhancedDuration
		decoded, err := hex.DecodeString(input)
		require.NoError(t, err)
		_, err = duration.ReadFrom(bytes.NewReader(decoded))
		require.NoError(t, err)
		require.Equal(t, expected, duration)
		var buf bytes.Buffer
		_, err = duration.WriteTo(&buf)
		require.NoError(t, err)
		require.Equal(t, decoded, buf.Bytes())
	}
}

func TestTime_ErrorHandler(t *testing.T) {
	var timestamp Time
	_, err := timestamp.ReadFrom(bytes.NewReader(nil))
	require.Error(t, err)
	var duration Duration
	_, err = duration.ReadFrom(bytes.NewReader(nil))
	require.Error(t, err)
	var enhancedDuration EnhancedDuration
	_, err = enhancedDuration.ReadFrom(bytes.NewReader(nil))
	require.Error(t, err)
}
