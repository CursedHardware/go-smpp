package sms

// Deliver see GSM 03.40, section 9.2.2.1 (35p)
type Deliver struct {
	SCAddress              SCAddress    `TP:"SC"`
	Flags                  DeliverFlags `DIR:"MT"`
	OriginatingAddress     Address      `TP:"OA"`
	ProtocolIdentifier     byte         `TP:"PID"`
	DataCoding             byte         `TP:"DCS"`
	ServiceCentreTimestamp Time         `TP:"SCTS"`
	UserData               []byte       `TP:"UD"`
}

// DeliverReport see GSM 03.40, section 9.2.2.1a (37p)
type DeliverReport struct {
	SCAddress          SCAddress          `TP:"SC"`
	Flags              Flags              `DIR:"MO"`
	ParameterIndicator ParameterIndicator `TP:"PI"`
	ProtocolIdentifier byte               `TP:"PID"`
	DataCoding         byte               `TP:"DCS"`
	UserData           []byte             `TP:"UD"`
}

type DeliverReportError struct {
	SCAddress    SCAddress    `TP:"SC"`
	Flags        Flags        `DIR:"MO"`
	FailureCause FailureCause `TP:"FCS"`
}

// Submit see GSM 03.40, section 9.2.2.2 (39p)
type Submit struct {
	SCAddress          SCAddress   `TP:"SC"`
	Flags              SubmitFlags `DIR:"MO"`
	MessageReference   byte        `TP:"MR"`
	DestinationAddress Address     `TP:"DA"`
	ProtocolIdentifier byte        `TP:"PID"`
	DataCoding         byte        `TP:"DCS"`
	ValidityPeriod     interface{} `TP:"VP"`
	UserData           []byte      `TP:"UD"`
}

// SubmitReport see GSM 03.40, section 9.2.2.2a (41p)
type SubmitReport struct {
	SCAddress              SCAddress          `TP:"SC"`
	Flags                  SubmitFlags        `DIR:"MT"`
	ParameterIndicator     ParameterIndicator `TP:"PI"`
	ServiceCentreTimestamp Time               `TP:"SCTS"`
	ProtocolIdentifier     byte               `TP:"PID"`
	DataCoding             byte               `TP:"DCS"`
	UserData               []byte             `TP:"UD"`
}

type SubmitReportError DeliverReportError

// StatusReport see GSM 03.40, section 9.2.2.3 (43p)
type StatusReport struct {
	SCAddress              SCAddress `TP:"SC"`
	Flags                  Flags     `DIR:"MT"`
	MessageReference       byte      `TP:"MR"`
	MoreMessagesToSend     bool      `TP:"MMS"`
	RecipientAddress       Address   `TP:"RA"`
	ServiceCentreTimestamp Time      `TP:"SCTS"`
	DischargeTime          Time      `TP:"DT"`
	Status                 byte      `TP:"ST"`
}

// Command see GSM 03.40, section 9.2.2.4 (45p)
type Command struct {
	SCAddress           SCAddress `TP:"SC"`
	Flags               Flags     `DIR:"MO"`
	MessageReference    byte      `TP:"MR"`
	StatusReportRequest bool      `TP:"SRR"`
	ProtocolIdentifier  byte      `TP:"PID"`
	CommandType         byte      `TP:"CT"`
	MessageNumber       byte      `TP:"MN"`
	DestinationAddress  Address   `TP:"DA"`
	CommandData         []byte    `TP:"CD"`
}
