package pdu

import . "github.com/M2MGateway/go-smpp/coding"

type Responsable interface {
	Resp() any
}

// AlertNotification see SMPP v5, section 4.1.3.1 (64p)
type AlertNotification struct {
	Header     Header
	SourceAddr Address
	ESMEAddr   Address
	Tags       Tags
}

// BindReceiver see SMPP v5, section 4.1.1.3 (58p)
type BindReceiver struct {
	Header       Header
	SystemID     string
	Password     string
	SystemType   string
	Version      InterfaceVersion
	AddressRange Address // see section 4.7.3.1
}

func (p *BindReceiver) Resp() any {
	return &BindReceiverResp{Header: Header{Sequence: p.Header.Sequence}, SystemID: p.SystemID}
}

// BindReceiverResp see SMPP v5, section 4.1.1.4 (59p)
type BindReceiverResp struct {
	Header   Header
	SystemID string
	Tags     Tags
}

// BindTransceiver see SMPP v5, section 4.1.1.5 (59p)
type BindTransceiver struct {
	Header       Header
	SystemID     string
	Password     string
	SystemType   string
	Version      InterfaceVersion
	AddressRange Address // see section 4.7.3.1
}

func (p *BindTransceiver) Resp() any {
	return &BindTransceiverResp{Header: Header{Sequence: p.Header.Sequence}, SystemID: p.SystemID}
}

// BindTransceiverResp see SMPP v5, section 4.1.1.6 (60p)
type BindTransceiverResp struct {
	Header   Header
	SystemID string
	Tags     Tags
}

// BindTransmitter see SMPP v5, section 4.1.1.1 (56p)
type BindTransmitter struct {
	Header       Header
	SystemID     string
	Password     string
	SystemType   string
	Version      InterfaceVersion
	AddressRange Address // see section 4.7.3.1
}

func (p *BindTransmitter) Resp() any {
	return &BindTransmitterResp{Header: Header{Sequence: p.Header.Sequence}, SystemID: p.SystemID}
}

// BindTransmitterResp see SMPP v5, section 4.1.1.2 (57p)
type BindTransmitterResp struct {
	Header   Header
	SystemID string
	Tags     Tags
}

// BroadcastSM see SMPP v5, section 4.4.1.1 (92p)
type BroadcastSM struct {
	Header               Header
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

func (p *BroadcastSM) Resp() any {
	return &BroadcastSMResp{Header: Header{Sequence: p.Header.Sequence}, MessageID: p.MessageID}
}

// BroadcastSMResp see SMPP v5, section 4.4.1.2 (96p)
type BroadcastSMResp struct {
	Header    Header
	MessageID string
	Tags      Tags
}

// CancelBroadcastSM see SMPP v5, section 4.6.2.1 (110p)
type CancelBroadcastSM struct {
	Header      Header
	ServiceType string
	MessageID   string
	SourceAddr  Address
	Tags        Tags
}

func (p *CancelBroadcastSM) Resp() any {
	return &CancelBroadcastSMResp{Header: Header{Sequence: p.Header.Sequence}}
}

// CancelBroadcastSMResp see SMPP v5, section 4.6.2.3 (112p)
type CancelBroadcastSMResp struct {
	Header Header
}

// CancelSM see SMPP v5, section 4.5.1.1 (100p)
type CancelSM struct {
	Header      Header
	ServiceType string
	MessageID   string
	SourceAddr  Address
	DestAddr    Address
}

func (p *CancelSM) Resp() any {
	return &CancelSMResp{Header: Header{Sequence: p.Header.Sequence}}
}

// CancelSMResp see SMPP v5, section 4.5.1.2 (101p)
type CancelSMResp struct {
	Header Header
}

// DataSM see SMPP v5, section 4.2.2.1 (69p)
type DataSM struct {
	Header             Header
	ServiceType        string
	SourceAddr         Address
	DestAddr           Address
	ESMClass           ESMClass
	RegisteredDelivery RegisteredDelivery
	DataCoding         DataCoding
	Tags               Tags
}

func (p *DataSM) Resp() any {
	return &DataSMResp{Header: Header{Sequence: p.Header.Sequence}}
}

// DataSMResp see SMPP v5, section 4.2.2.2 (70p)
type DataSMResp struct {
	Header    Header
	MessageID string
	Tags      Tags
}

// DeliverSM see SMPP v5, section 4.3.1.1 (85p)
type DeliverSM struct {
	Header               Header
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

func (p *DeliverSM) Resp() any {
	return &DeliverSMResp{Header: Header{Sequence: p.Header.Sequence}}
}

// DeliverSMResp see SMPP v5, section 4.3.1.1 (87p)
type DeliverSMResp struct {
	Header    Header
	MessageID string
	Tags      Tags
}

// EnquireLink see SMPP v5, section 4.1.2.1 (63p)
type EnquireLink struct {
	Header Header
	Tags   Tags
}

func (p *EnquireLink) Resp() any {
	return &EnquireLinkResp{Header: Header{Sequence: p.Header.Sequence}}
}

// EnquireLinkResp see SMPP v5, section 4.1.2.2 (63p)
type EnquireLinkResp struct {
	Header Header
}

// GenericNACK see SMPP v5, section 4.1.4.1 (65p)
type GenericNACK struct {
	Header Header
	Tags   Tags
}

// Outbind see SMPP v5, section 4.1.1.7 (61p)
type Outbind struct {
	Header   Header
	SystemID string
	Password string
}

// QueryBroadcastSM see SMPP v5, section 4.6.1.1 (107p)
type QueryBroadcastSM struct {
	Header     Header
	MessageID  string
	SourceAddr Address
	Tags       Tags
}

func (p *QueryBroadcastSM) Resp() any {
	return &QueryBroadcastSMResp{Header: Header{Sequence: p.Header.Sequence}, MessageID: p.MessageID}
}

// QueryBroadcastSMResp see SMPP v5, section 4.6.1.3 (108p)
type QueryBroadcastSMResp struct {
	Header    Header
	MessageID string
	Tags      Tags
}

// QuerySM see SMPP v5, section 4.5.2.1 (101p)
type QuerySM struct {
	Header     Header
	MessageID  string
	SourceAddr Address
}

func (p *QuerySM) Resp() any {
	return &QuerySMResp{Header: Header{Sequence: p.Header.Sequence}}
}

// QuerySMResp see SMPP v5, section 4.5.2.2 (103p)
type QuerySMResp struct {
	Header       Header
	MessageID    string
	FinalDate    string
	MessageState MessageState
	ErrorCode    CommandStatus
}

// ReplaceSM see SMPP v5, section 4.5.3.1 (104p)
type ReplaceSM struct {
	Header               Header
	MessageID            string
	SourceAddr           Address
	ScheduleDeliveryTime string
	ValidityPeriod       string
	RegisteredDelivery   RegisteredDelivery
	Message              ShortMessage
	Tags                 Tags
}

func (p *ReplaceSM) Resp() any {
	return &ReplaceSMResp{Header: Header{Sequence: p.Header.Sequence}}
}

// ReplaceSMResp see SMPP v5, section 4.5.3.2 (106p)
type ReplaceSMResp struct {
	Header Header
}

// SubmitMulti see SMPP v5, section 4.2.3.1 (71p)
type SubmitMulti struct {
	Header               Header
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

func (p *SubmitMulti) Resp() any {
	return &SubmitMultiResp{Header: Header{Sequence: p.Header.Sequence}}
}

// SubmitMultiResp see SMPP v5, section 4.2.3.2 (74p)
type SubmitMultiResp struct {
	Header           Header
	MessageID        string
	UnsuccessfulSMEs UnsuccessfulRecords
	Tags             Tags
}

// SubmitSM see SMPP v5, section 4.2.1.1 (66p)
type SubmitSM struct {
	Header               Header
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

func (p *SubmitSM) Resp() any {
	return &SubmitSMResp{Header: Header{Sequence: p.Header.Sequence}}
}

// SubmitSMResp see SMPP v5, section 4.2.1.2 (68p)
type SubmitSMResp struct {
	Header    Header
	MessageID string
}

// Unbind see SMPP v5, section 4.1.1.8 (61p)
type Unbind struct {
	Header Header
}

func (p *Unbind) Resp() any {
	return &UnbindResp{Header: Header{Sequence: p.Header.Sequence}}
}

// UnbindResp see SMPP v5, section 4.1.1.9 (62p)
type UnbindResp struct {
	Header Header
}
