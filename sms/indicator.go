package sms

import "fmt"

type ParameterIndicator struct {
	ProtocolIdentifier bool
	DataCoding         bool
	UserData           bool
}

func (p *ParameterIndicator) Has(abbr string) bool {
	switch abbr {
	case "PID":
		return p.ProtocolIdentifier
	case "DCS":
		return p.DataCoding
	case "UD":
		return p.UserData
	}
	return false
}

func (p *ParameterIndicator) Set(abbr string) {
	switch abbr {
	case "PID":
		p.ProtocolIdentifier = true
	case "DCS":
		p.DataCoding = true
	case "UD":
		p.UserData = true
	}
}

func (p *ParameterIndicator) WriteByte(c byte) error {
	return unmarshalFlags(c, p)
}

func (p *ParameterIndicator) ReadByte() (byte, error) {
	return marshalFlags(p)
}

type Flags struct {
	MessageType MessageType
}

func (p *Flags) setDirection(direction Direction) {
	p.MessageType.Set(p.MessageType.Type(), direction)
}

func (p *Flags) WriteByte(c byte) error {
	return unmarshalFlags(c, p)
}

func (p *Flags) ReadByte() (byte, error) {
	return marshalFlags(p)
}

type DeliverFlags struct {
	MessageType            MessageType
	MoreMessagesToSend     bool
	ReplyPath              bool
	UDHIndicator           bool
	StatusReportIndication bool
}

func (p *DeliverFlags) setDirection(direction Direction) {
	p.MessageType.Set(p.MessageType.Type(), direction)
}

func (p *DeliverFlags) WriteByte(c byte) error {
	return unmarshalFlags(c, p)
}

func (p *DeliverFlags) ReadByte() (byte, error) {
	return marshalFlags(p)
}

type SubmitFlags struct {
	MessageType             MessageType
	RejectDuplicates        bool
	ValidityPeriodFormat    byte
	ReplyPath               bool
	UserDataHeaderIndicator bool
	StatusReportRequest     bool
}

func (p *SubmitFlags) setDirection(direction Direction) {
	p.MessageType.Set(p.MessageType.Type(), direction)
}

func (p *SubmitFlags) WriteByte(c byte) error {
	return unmarshalFlags(c, p)
}

func (p *SubmitFlags) ReadByte() (byte, error) {
	return marshalFlags(p)
}

// FailureCause see GSM 03.40, section 9.2.3.22 (54p)
type FailureCause byte

//goland:noinspection SpellCheckingInspection
var failureCauseErrors = map[FailureCause]string{
	0x00: "Telematic interworking not supported",
	0x01: "Short message Type 0 not supported",
	0x02: "Cannot replace short message",
	0x0F: "Unspecified TP-PID error",
	0x10: "Data coding schema (alphabet not supported)",
	0x11: "Message class not supported",
	0x1F: "Unspecified TP-DCS error",
	0x20: "Command cannot be actioned",
	0x21: "Command unsupported",
	0x2F: "Unspecified TP-Command error",
	0x30: "TPDU not supported",
	0x40: "SC busy",
	0x41: "No SC subscription",
	0x42: "SC system failure",
	0x43: "Invalid SME address",
	0x44: "Destination SME barred",
	0x45: "SME Rejected-Duplicate",
	0x50: "SIM SMS storage full",
	0x51: "No SMS storage capability in SIM",
	0x52: "Error in MS",
	0x53: "Memory Capacity Exceeded",
	0x7F: "Unspecified error cause",
}

//goland:noinspection SpellCheckingInspection
func (f FailureCause) Error() string {
	if message, ok := failureCauseErrors[f-0x80]; ok {
		return message
	}
	return fmt.Sprintf("%X", byte(f))
}
