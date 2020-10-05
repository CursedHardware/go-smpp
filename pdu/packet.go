package pdu

import . "github.com/NiceLabs/go-smpp/coding"

type Responsable interface {
	Resp() interface{}
}

// AlertNotification see SMPP v5, section 4.1.3.1 (64p)
type AlertNotification struct {
	Header     Header `id:"00000102"`
	SourceAddr Address
	ESMEAddr   Address
	Tags       Tags
}

// BindReceiver see SMPP v5, section 4.1.1.3 (58p)
type BindReceiver struct {
	Header       Header `id:"00000001"`
	SystemID     string
	Password     string
	SystemType   string
	Version      InterfaceVersion
	AddressRange Address // see section 4.7.3.1
}

func (p *BindReceiver) Resp() interface{} {
	return &BindReceiverResp{Header: Header{Sequence: p.Header.Sequence}, SystemID: p.SystemID}
}

// BindReceiverResp see SMPP v5, section 4.1.1.4 (59p)
type BindReceiverResp struct {
	Header   Header `id:"80000001"`
	SystemID string
	Tags     Tags
}

// BindTransceiver see SMPP v5, section 4.1.1.5 (59p)
type BindTransceiver struct {
	Header       Header `id:"00000009"`
	SystemID     string
	Password     string
	SystemType   string
	Version      InterfaceVersion
	AddressRange Address // see section 4.7.3.1
}

func (p *BindTransceiver) Resp() interface{} {
	return &BindTransceiverResp{Header: Header{Sequence: p.Header.Sequence}, SystemID: p.SystemID}
}

// BindTransceiverResp see SMPP v5, section 4.1.1.6 (60p)
type BindTransceiverResp struct {
	Header   Header `id:"80000009"`
	SystemID string
	Tags     Tags
}

// BindTransmitter see SMPP v5, section 4.1.1.1 (56p)
type BindTransmitter struct {
	Header       Header `id:"00000002"`
	SystemID     string
	Password     string
	SystemType   string
	Version      InterfaceVersion
	AddressRange Address // see section 4.7.3.1
}

func (p *BindTransmitter) Resp() interface{} {
	return &BindTransmitterResp{Header: Header{Sequence: p.Header.Sequence}, SystemID: p.SystemID}
}

// BindTransmitter see SMPP v5, section 4.1.1.2 (57p)
type BindTransmitterResp struct {
	Header   Header `id:"80000002"`
	SystemID string
	Tags     Tags
}

// BroadcastSM see SMPP v5, section 4.4.1.1 (92p)
type BroadcastSM struct {
	Header               Header `id:"00000112"`
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

func (p *BroadcastSM) Resp() interface{} {
	return &BroadcastSMResp{Header: Header{Sequence: p.Header.Sequence}, MessageID: p.MessageID}
}

// BroadcastSMResp see SMPP v5, section 4.4.1.2 (96p)
type BroadcastSMResp struct {
	Header    Header `id:"80000112"`
	MessageID string
	Tags      Tags
}

// CancelBroadcastSM see SMPP v5, section 4.6.2.1 (110p)
type CancelBroadcastSM struct {
	Header      Header `id:"00000113"`
	ServiceType string
	MessageID   string
	SourceAddr  Address
	Tags        Tags
}

func (p *CancelBroadcastSM) Resp() interface{} {
	return &CancelBroadcastSMResp{Header: Header{Sequence: p.Header.Sequence}}
}

// CancelBroadcastSMResp see SMPP v5, section 4.6.2.3 (112p)
type CancelBroadcastSMResp struct {
	Header Header `id:"80000113"`
}

// CancelSM see SMPP v5, section 4.5.1.1 (100p)
type CancelSM struct {
	Header      Header `id:"00000008"`
	ServiceType string
	MessageID   string
	SourceAddr  Address
	DestAddr    Address
}

func (p *CancelSM) Resp() interface{} {
	return &CancelSMResp{Header: Header{Sequence: p.Header.Sequence}}
}

// CancelSM see SMPP v5, section 4.5.1.2 (101p)
type CancelSMResp struct {
	Header Header `id:"80000008"`
}

// DataSM see SMPP v5, section 4.2.2.1 (69p)
type DataSM struct {
	Header             Header `id:"00000103"`
	ServiceType        string
	SourceAddr         Address
	DestAddr           Address
	ESMClass           ESMClass
	RegisteredDelivery RegisteredDelivery
	DataCoding         DataCoding
	Tags               Tags
}

func (p *DataSM) Resp() interface{} {
	return &DataSMResp{Header: Header{Sequence: p.Header.Sequence}}
}

// DataSM see SMPP v5, section 4.2.2.2 (70p)
type DataSMResp struct {
	Header    Header `id:"80000103"`
	MessageID string
	Tags      Tags
}

// DeliverSM see SMPP v5, section 4.3.1.1 (85p)
type DeliverSM struct {
	Header               Header `id:"00000005"`
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

func (p *DeliverSM) Resp() interface{} {
	return &DeliverSMResp{Header: Header{Sequence: p.Header.Sequence}}
}

// DeliverSMResp see SMPP v5, section 4.3.1.1 (87p)
type DeliverSMResp struct {
	Header    Header `id:"80000005"`
	MessageID string
	Tags      Tags
}

// EnquireLink see SMPP v5, section 4.1.2.1 (63p)
type EnquireLink struct {
	Header Header `id:"00000015"`
	Tags   Tags
}

func (p *EnquireLink) Resp() interface{} {
	return &EnquireLinkResp{Header: Header{Sequence: p.Header.Sequence}}
}

// EnquireLink see SMPP v5, section 4.1.2.2 (63p)
type EnquireLinkResp struct {
	Header Header `id:"80000015"`
}

// GenericNACK see SMPP v5, section 4.1.4.1 (65p)
type GenericNACK struct {
	Header Header `id:"80000000"`
	Tags   Tags
}

// Outbind see SMPP v5, section 4.1.1.7 (61p)
type Outbind struct {
	Header   Header `id:"0000000B"`
	SystemID string
	Password string
}

// QueryBroadcastSM see SMPP v5, section 4.6.1.1 (107p)
type QueryBroadcastSM struct {
	Header     Header `id:"00000111"`
	MessageID  string
	SourceAddr Address
	Tags       Tags
}

func (p *QueryBroadcastSM) Resp() interface{} {
	return &QueryBroadcastSMResp{Header: Header{Sequence: p.Header.Sequence}, MessageID: p.MessageID}
}

// QueryBroadcastSMResp see SMPP v5, section 4.6.1.3 (108p)
type QueryBroadcastSMResp struct {
	Header    Header `id:"80000111"`
	MessageID string
	Tags      Tags
}

// QuerySM see SMPP v5, section 4.5.2.1 (101p)
type QuerySM struct {
	Header     Header `id:"00000003"`
	MessageID  string
	SourceAddr Address
}

func (p *QuerySM) Resp() interface{} {
	return &QuerySMResp{Header: Header{Sequence: p.Header.Sequence}}
}

// QuerySMResp see SMPP v5, section 4.5.2.2 (103p)
type QuerySMResp struct {
	Header       Header `id:"80000003"`
	MessageID    string
	FinalDate    string
	MessageState byte // section 4.7.15, page 127
	ErrorCode    CommandStatus
}

// ReplaceSM see SMPP v5, section 4.5.3.1 (104p)
type ReplaceSM struct {
	Header               Header `id:"00000007"`
	MessageID            string
	SourceAddr           Address
	ScheduleDeliveryTime string
	ValidityPeriod       string
	RegisteredDelivery   RegisteredDelivery
	Message              ShortMessage
	Tags                 Tags
}

func (p *ReplaceSM) Resp() interface{} {
	return &ReplaceSMResp{Header: Header{Sequence: p.Header.Sequence}}
}

// ReplaceSMResp see SMPP v5, section 4.5.3.2 (106p)
type ReplaceSMResp struct {
	Header Header `id:"80000007"`
}

// SubmitMulti see SMPP v5, section 4.2.3.1 (71p)
type SubmitMulti struct {
	Header               Header `id:"00000021"`
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

func (p *SubmitMulti) Resp() interface{} {
	return &SubmitMultiResp{Header: Header{Sequence: p.Header.Sequence}}
}

// SubmitMulti see SMPP v5, section 4.2.3.2 (74p)
type SubmitMultiResp struct {
	Header           Header `id:"80000021"`
	MessageID        string
	UnsuccessfulSMEs UnsuccessfulRecords
	Tags             Tags
}

// SubmitSM see SMPP v5, section 4.2.1.1 (66p)
type SubmitSM struct {
	Header               Header `id:"00000004"`
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

func (p *SubmitSM) Resp() interface{} {
	return &SubmitSMResp{Header: Header{Sequence: p.Header.Sequence}}
}

// SubmitSMResp see SMPP v5, section 4.2.1.2 (68p)
type SubmitSMResp struct {
	Header    Header `id:"80000004"`
	MessageID string
}

// Unbind see SMPP v5, section 4.1.1.8 (61p)
type Unbind struct {
	Header Header `id:"00000006"`
}

func (p *Unbind) Resp() interface{} {
	return &UnbindResp{Header: Header{Sequence: p.Header.Sequence}}
}

// UnbindResp see SMPP v5, section 4.1.1.9 (62p)
type UnbindResp struct {
	Header Header `id:"80000006"`
}
