package sms

type Flags struct {
	MessageType MessageType
}

type DeliverFlags struct {
	MessageType            MessageType
	MoreMessagesToSend     bool
	ReplyPath              bool
	UDHIndicator           bool
	StatusReportIndication bool
}

type SubmitFlags struct {
	MessageType             MessageType
	RejectDuplicates        bool
	ValidityPeriodFormat    byte
	ReplyPath               bool
	UserDataHeaderIndicator bool
	StatusReportRequest     bool
}
