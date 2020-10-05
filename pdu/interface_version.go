package pdu

// InterfaceVersion see SMPP v5, section 4.7.13 (126p)
type InterfaceVersion byte

const (
	SMPPVersion33 InterfaceVersion = 0x33
	SMPPVersion34 InterfaceVersion = 0x34
	SMPPVersion50 InterfaceVersion = 0x50
)
