package sms

// Deliver see GSM 03.40, section 9.2.2.1 (35p)
type Deliver struct {
	Flags                  DeliverFlags `DIR:"MT"`
	OriginatingAddress     Address      `TP:"OA"`
	ProtocolIdentifier     byte         `TP:"PI"`
	DataCoding             byte         `TP:"DCS"`
	ServiceCentreTimestamp Time         `TP:"SCTS"`
	UserData               []byte       `TP:"UD"`
}

// DeliverReport see GSM 03.40, section 9.2.2.1a (37p)
type DeliverReport struct {
	Flags              Flags              `DIR:"MO"`
	ParameterIndicator ParameterIndicator `TP:"PI"`
	ProtocolIdentifier byte               `TP:"PID"`
	DataCoding         byte               `TP:"DCS"`
	UserData           []byte             `TP:"UD"`
}

type DeliverReportError struct {
	Flags        Flags `DIR:"MO"`
	FailureCause byte  `TP:"FCS"`
}

// Submit see GSM 03.40, section 9.2.2.2 (39p)
type Submit struct {
	Flags              SubmitFlags    `DIR:"MO"`
	MessageReference   byte           `TP:"MR"`
	DestinationAddress Address        `TP:"DA"`
	ProtocolIdentifier byte           `TP:"PI"`
	DataCoding         byte           `TP:"DCS"`
	ValidityPeriod     ValidityPeriod `TP:"VP"`
	UserData           []byte         `TP:"UD"`
}

// SubmitReport see GSM 03.40, section 9.2.2.2a (41p)
type SubmitReport struct {
	Flags                  SubmitFlags        `DIR:"MT"`
	ParameterIndicator     ParameterIndicator `TP:"PI"`
	ServiceCentreTimestamp Time               `TP:"SCTS"`
	ProtocolIdentifier     byte               `TP:"PID"`
	DataCoding             byte               `TP:"DCS"`
	UserData               []byte             `TP:"UD"`
}

type SubmitReportError struct {
	Flags        Flags `DIR:"MT"`
	FailureCause byte  `TP:"FCS"`
}

// StatusReport see GSM 03.40, section 9.2.2.3 (43p)
type StatusReport struct {
	Flags                  Flags   `DIR:"MT"`
	MessageReference       byte    `TP:"MR"`
	MoreMessagesToSend     bool    `TP:"MMS"`
	RecipientAddress       Address `TP:"RA"`
	ServiceCentreTimestamp Time    `TP:"SCTS"`
	DischargeTime          Time    `TP:"DT"`
	Status                 byte    `TP:"ST"`
}

// Command see GSM 03.40, section 9.2.2.4 (45p)
type Command struct {
	Flags               Flags   `DIR:"MO"`
	MessageReference    byte    `TP:"MR"`
	StatusReportRequest bool    `TP:"SRR"`
	ProtocolIdentifier  byte    `TP:"PID"`
	CommandType         byte    `TP:"CT"`
	MessageNumber       byte    `TP:"MN"`
	DestinationAddress  Address `TP:"DA"`
	CommandData         []byte  `TP:"CD"`
}
