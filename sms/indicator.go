package sms

import "fmt"

type ParameterIndicator struct {
	ProtocolIdentifier bool
	DataCoding         bool
	UserData           bool
}

func (p ParameterIndicator) get(abbr string) *bool {
	switch abbr {
	case "PID":
		return &p.ProtocolIdentifier
	case "DCS":
		return &p.DataCoding
	case "UD":
		return &p.UserData
	}
	return nil
}

func (p ParameterIndicator) Has(abbr string) bool {
	v := p.get(abbr)
	return v != nil && *v
}

func (p *ParameterIndicator) Set(abbr string) {
	v := p.get(abbr)
	*v = true
}

func (p *ParameterIndicator) WriteByte(c byte) error {
	return unmarshalFlags(c, p)
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

// FailureCause see GSM 03.40, section 9.2.3.22 (54p)
type FailureCause byte

//goland:noinspection SpellCheckingInspection
var failureCauseErrors = map[FailureCause]string{
	0x80: "Telematic interworking not supported",
	0x81: "Short message Type 0 not supported",
	0x82: "Cannot replace short message",
	0x8F: "Unspecified TP-PID error",
	0x90: "Data coding schema (alphabet not supported)",
	0x91: "Message class not supported",
	0x9F: "Unspecified TP-DCS error",
	0xA0: "Command cannot be actioned",
	0xA1: "Command unsupported",
	0xAF: "Unspecified TP-Command error",
	0xB0: "TPDU not supported",
	0xC0: "SC busy",
	0xC1: "No SC subscription",
	0xC2: "SC system failure",
	0xC3: "Invalid SME address",
	0xC4: "Destination SME barred",
	0xC5: "SME Rejected-Duplicate",
	0xD0: "SIM SMS storage full",
	0xD1: "No SMS storage capability in SIM",
	0xD2: "Error in MS",
	0xD3: "Memory Capacity Exceeded",
	0xFF: "Unspecified error cause",
}

//goland:noinspection SpellCheckingInspection
func (f FailureCause) Error() string {
	if message, ok := failureCauseErrors[f]; ok {
		return message
	}
	return fmt.Sprintf("%X", byte(f))
}
