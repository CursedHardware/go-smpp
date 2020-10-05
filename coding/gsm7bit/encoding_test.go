package gsm7bit

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

//goland:noinspection SpellCheckingInspection
var mapping = map[string]string{
	"":                                "",
	"1":                               "31",
	"12":                              "3119",
	"123":                             "31D90C",
	"1234":                            "31D98C06",
	"12345":                           "31D98C5603",
	"123456":                          "31D98C56B301",
	"1234567":                         "31D98C56B3DD1A",
	"12345678":                        "31D98C56B3DD70",
	"123456789":                       "31D98C56B3DD7039",
	"12345[6]":                        "31D98C56DBF06C1B1F",
	"^{}\\[~]|€":                      "1bca06b5496d5e1bdea6b7f16d809b32",
	"of the printing and typesetting": "6F33888E2E83E0F2B49B9E769F4161371944CFC3CBF3329D9E769F1B",
	"industry. Lorem Ipsum has been":  "6937B93EA7CBF32E10F32D2FB74149F8BCDE06A1C37390B85C7603",
	"the industry's standard dummy":   "747419947693EB73BA3C7F9A83E6F4B09B1C969341E47ABB9D07",
}

var invalidEncoder = [][]byte{
	{0xFF},
}

var invalidDecoder = [][]byte{
	{0x1B},
	{0x1B, 0x80},
}

func TestGSM7Encoding(t *testing.T) {
	encoder := Packed.NewEncoder()
	decoder := Packed.NewDecoder()
	for decodedText, encodedHex := range mapping {
		decoded, err := hex.DecodeString(encodedHex)
		require.NoError(t, err)

		encoded, err := encoder.Bytes([]byte(decodedText))
		require.NoError(t, err)
		require.Equal(t, decoded, encoded, hex.EncodeToString(encoded))

		decoded, err = decoder.Bytes(encoded)
		require.NoError(t, err)
		require.Equal(t, decodedText, string(decoded), hex.EncodeToString(encoded))
	}
	for _, encoded := range invalidDecoder {
		_, _ = decoder.Bytes(encoded)
	}
}

func TestGSM7Encoding_invalid(t *testing.T) {
	encoder := Packed.NewEncoder()
	decoder := Packed.NewDecoder()
	for _, input := range invalidEncoder {
		_, err := encoder.Bytes(input)
		require.Error(t, err)
	}
	for _, encoded := range invalidDecoder {
		_, err := decoder.Bytes(encoded)
		require.Error(t, err)
	}
}

func TestGSM7Encoding_edge_case(t *testing.T) {
	encoder := Packed.NewEncoder()
	_, _ = encoder.Bytes([]byte("1234567\r"))
}
