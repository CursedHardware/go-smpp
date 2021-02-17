package pdu

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterfaceVersion(t *testing.T) {
	samples := map[InterfaceVersion][]byte{
		SMPPVersion33: []byte(`"3.3"`),
		SMPPVersion34: []byte(`"3.4"`),
		SMPPVersion50: []byte(`"5.0"`),
	}
	var err error
	var version InterfaceVersion
	for expected, expectedEncoded := range samples {
		err = json.Unmarshal(expectedEncoded, &version)
		assert.NoError(t, err)
		assert.Equal(t, expected, version)

		encoded, err := json.Marshal(expected)
		assert.NoError(t, err)
		assert.Equal(t, expectedEncoded, encoded)
	}

}
