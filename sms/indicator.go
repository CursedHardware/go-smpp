package sms

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
