package pdu

import . "github.com/M2MGateway/go-smpp/coding"

type Packet interface {
	getCommandStatus() CommandStatus
	getSequence() int32
	setSequence(int32)
}

// Responsable Response-able PDU interface
type Responsable interface {
	Packet
	Resp() Packet
}

// AlertNotification see SMPP v5, section 4.1.3.1 (64p)
type AlertNotification struct {
	Header     `id:"00000102"`
	SourceAddr Address
	ESMEAddr   Address
	Tags       Tags
}

// BindReceiver see SMPP v5, section 4.1.1.3 (58p)
type BindReceiver struct {
	Header       `id:"00000001"`
	SystemID     string
	Password     string
	SystemType   string
	Version      InterfaceVersion
	AddressRange Address // see section 4.7.3.1
}

// BindReceiverResp see SMPP v5, section 4.1.1.4 (59p)
type BindReceiverResp struct {
	Header   `id:"80000001"`
	SystemID string
	Tags     Tags
}

// BindTransceiver see SMPP v5, section 4.1.1.5 (59p)
type BindTransceiver struct {
	Header       `id:"00000009"`
	SystemID     string
	Password     string
	SystemType   string
	Version      InterfaceVersion
	AddressRange Address // see section 4.7.3.1
}

// BindTransceiverResp see SMPP v5, section 4.1.1.6 (60p)
type BindTransceiverResp struct {
	Header   `id:"80000009"`
	SystemID string
	Tags     Tags
}

// BindTransmitter see SMPP v5, section 4.1.1.1 (56p)
type BindTransmitter struct {
	Header       `id:"00000002"`
	SystemID     string
	Password     string
	SystemType   string
	Version      InterfaceVersion
	AddressRange Address // see section 4.7.3.1
}

// BindTransmitterResp see SMPP v5, section 4.1.1.2 (57p)
type BindTransmitterResp struct {
	Header   `id:"80000002"`
	SystemID string
	Tags     Tags
}

// BroadcastSM see SMPP v5, section 4.4.1.1 (92p)
type BroadcastSM struct {
	Header               `id:"00000112"`
	ServiceType          string
	SourceAddr           Address
	MessageID            string
	PriorityFlag         byte
	ScheduleDeliveryTime string
	ValidityPeriod       string
	ReplaceIfPresent     bool
	DataCoding           DataCoding
	DefaultMessageID     byte
	Tags                 Tags
}

// BroadcastSMResp see SMPP v5, section 4.4.1.2 (96p)
type BroadcastSMResp struct {
	Header    `id:"80000112"`
	MessageID string
	Tags      Tags
}

// CancelBroadcastSM see SMPP v5, section 4.6.2.1 (110p)
type CancelBroadcastSM struct {
	Header      `id:"00000113"`
	ServiceType string
	MessageID   string
	SourceAddr  Address
	Tags        Tags
}

// CancelBroadcastSMResp see SMPP v5, section 4.6.2.3 (112p)
type CancelBroadcastSMResp struct {
	Header `id:"80000113"`
}

// CancelSM see SMPP v5, section 4.5.1.1 (100p)
type CancelSM struct {
	Header      `id:"00000008"`
	ServiceType string
	MessageID   string
	SourceAddr  Address
	DestAddr    Address
}

// CancelSMResp CancelSM see SMPP v5, section 4.5.1.2 (101p)
type CancelSMResp struct {
	Header `id:"80000008"`
}

// DataSM see SMPP v5, section 4.2.2.1 (69p)
type DataSM struct {
	Header             `id:"00000103"`
	ServiceType        string
	SourceAddr         Address
	DestAddr           Address
	ESMClass           ESMClass
	RegisteredDelivery RegisteredDelivery
	DataCoding         DataCoding
	Tags               Tags
}

// DataSMResp see SMPP v5, section 4.2.2.2 (70p)
type DataSMResp struct {
	Header    `id:"80000103"`
	MessageID string
	Tags      Tags
}

// DeliverSM see SMPP v5, section 4.3.1.1 (85p)
type DeliverSM struct {
	Header               `id:"00000005"`
	ServiceType          string
	SourceAddr           Address
	DestAddr             Address
	ESMClass             ESMClass
	ProtocolID           byte
	PriorityFlag         byte
	ScheduleDeliveryTime string
	ValidityPeriod       string
	RegisteredDelivery   RegisteredDelivery
	ReplaceIfPresent     bool
	Message              ShortMessage
	Tags                 Tags
}

// DeliverSMResp see SMPP v5, section 4.3.1.1 (87p)
type DeliverSMResp struct {
	Header    `id:"80000005"`
	MessageID string
	Tags      Tags
}

// EnquireLink see SMPP v5, section 4.1.2.1 (63p)
type EnquireLink struct {
	Header `id:"00000015"`
	Tags   Tags
}

// EnquireLinkResp see SMPP v5, section 4.1.2.2 (63p)
type EnquireLinkResp struct {
	Header `id:"80000015"`
}

// GenericNACK see SMPP v5, section 4.1.4.1 (65p)
type GenericNACK struct {
	Header `id:"80000000"`
	Tags   Tags
}

// Outbind see SMPP v5, section 4.1.1.7 (61p)
type Outbind struct {
	Header   `id:"0000000B"`
	SystemID string
	Password string
}

// QueryBroadcastSM see SMPP v5, section 4.6.1.1 (107p)
type QueryBroadcastSM struct {
	Header     `id:"00000111"`
	MessageID  string
	SourceAddr Address
	Tags       Tags
}

// QueryBroadcastSMResp see SMPP v5, section 4.6.1.3 (108p)
type QueryBroadcastSMResp struct {
	Header    `id:"80000111"`
	MessageID string
	Tags      Tags
}

// QuerySM see SMPP v5, section 4.5.2.1 (101p)
type QuerySM struct {
	Header     `id:"00000003"`
	MessageID  string
	SourceAddr Address
}

// QuerySMResp see SMPP v5, section 4.5.2.2 (103p)
type QuerySMResp struct {
	Header       `id:"80000003"`
	MessageID    string
	FinalDate    string
	MessageState MessageState
	ErrorCode    CommandStatus
}

// ReplaceSM see SMPP v5, section 4.5.3.1 (104p)
type ReplaceSM struct {
	Header               `id:"00000007"`
	MessageID            string
	SourceAddr           Address
	ScheduleDeliveryTime string
	ValidityPeriod       string
	RegisteredDelivery   RegisteredDelivery
	Message              ShortMessage
	Tags                 Tags
}

// ReplaceSMResp see SMPP v5, section 4.5.3.2 (106p)
type ReplaceSMResp struct {
	Header `id:"80000007"`
}

// SubmitMulti see SMPP v5, section 4.2.3.1 (71p)
type SubmitMulti struct {
	Header               `id:"00000021"`
	ServiceType          string
	SourceAddr           Address
	DestAddrList         DestinationAddresses
	ESMClass             ESMClass
	ProtocolID           byte
	PriorityFlag         byte
	ScheduleDeliveryTime string
	ValidityPeriod       string
	RegisteredDelivery   RegisteredDelivery
	ReplaceIfPresent     bool
	Message              ShortMessage
	Tags                 Tags
}

// SubmitMultiResp see SMPP v5, section 4.2.3.2 (74p)
type SubmitMultiResp struct {
	Header           `id:"80000021"`
	MessageID        string
	UnsuccessfulSMEs UnsuccessfulRecords
	Tags             Tags
}

// SubmitSM see SMPP v5, section 4.2.1.1 (66p)
type SubmitSM struct {
	Header               `id:"00000004"`
	ServiceType          string
	SourceAddr           Address
	DestAddr             Address
	ESMClass             ESMClass
	ProtocolID           byte
	PriorityFlag         byte
	ScheduleDeliveryTime string
	ValidityPeriod       string
	RegisteredDelivery   RegisteredDelivery
	ReplaceIfPresent     bool
	Message              ShortMessage
	Tags                 Tags
}

// SubmitSMResp see SMPP v5, section 4.2.1.2 (68p)
type SubmitSMResp struct {
	Header    `id:"80000004"`
	MessageID string
}

// Unbind see SMPP v5, section 4.1.1.8 (61p)
type Unbind struct {
	Header `id:"00000006"`
}

// UnbindResp see SMPP v5, section 4.1.1.9 (62p)
type UnbindResp struct {
	Header `id:"80000006"`
}
