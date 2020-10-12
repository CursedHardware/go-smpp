package sms

type Error struct {
	Flags        Flags
	FailureCause byte
	UserData     []byte
}

// Deliver see GSM 03.40, section 9.2.2.1 (35p)
type Deliver struct {
	Flags                  DeliverFlags
	OriginatingAddress     Address
	ProtocolIdentifier     byte
	DataCoding             byte
	ServiceCentreTimestamp Time
	UserData               []byte
}

// DeliverReport see GSM 03.40, section 9.2.2.1a (37p)
type DeliverReport struct {
	Flags              Flags
	ParameterIndicator byte
	ProtocolIdentifier byte
	DataCoding         byte
	UserData           []byte
}

// Submit see GSM 03.40, section 9.2.2.2 (39p)
type Submit struct {
	Flags              SubmitFlags
	MessageReference   byte
	DestinationAddress Address
	ProtocolIdentifier byte
	DataCoding         byte
	ValidityPeriod     ValidityPeriod
	UserData           []byte
}

// SubmitReport see GSM 03.40, section 9.2.2.2a (41p)
type SubmitReport struct {
	Flags                  SubmitFlags
	ParameterIndicator     byte
	ServiceCentreTimestamp Time
	ProtocolIdentifier     byte
	DataCoding             byte
	UserData               []byte
}

// StatusReport see GSM 03.40, section 9.2.2.3 (43p)
type StatusReport struct {
	Flags                  Flags
	MessageReference       byte
	MoreMessagesToSend     bool
	RecipientAddress       Address
	ServiceCentreTimestamp Time
	DischargeTime          Time
	Status                 byte
}

// Command see GSM 03.40, section 9.2.2.4 (45p)
type Command struct {
	Flags               Flags
	MessageReference    byte
	StatusReportRequest bool
	ProtocolIdentifier  byte
	CommandType         byte
	MessageNumber       byte
	DestinationAddress  Address
	CommandData         []byte
}
