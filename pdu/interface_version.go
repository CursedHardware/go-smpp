package pdu

import (
	"encoding/json"
	"fmt"
)

// InterfaceVersion see SMPP v5, section 4.7.13 (126p)
type InterfaceVersion byte

const (
	SMPPVersion33 InterfaceVersion = 0x33
	SMPPVersion34 InterfaceVersion = 0x34
	SMPPVersion50 InterfaceVersion = 0x50
)

func (v InterfaceVersion) String() string {
	major := (v >> 4) & 0b1111
	minor := v & 0b1111
	return fmt.Sprintf("v%d.%d", major, minor)
}

func (v InterfaceVersion) MarshalJSON() (data []byte, err error) {
	return json.Marshal(v.String())
}

func (v *InterfaceVersion) UnmarshalJSON(data []byte) (err error) {
	var value string
	var major, minor InterfaceVersion
	if err = json.Unmarshal(data, &value); err != nil {
		return
	}
	if _, err = fmt.Sscanf(value, "v%d.%d", &major, &minor); err != nil {
		return
	}
	*v = ((major & 0b1111) << 4) | (minor & 0b1111)
	return
}
