package pdu

func (p *BindReceiver) Resp() Packet {
	return &BindReceiverResp{
		Header:   p.makeHeader(),
		SystemID: p.SystemID,
	}
}

func (p *BindTransceiver) Resp() Packet {
	return &BindTransceiverResp{
		Header:   p.makeHeader(),
		SystemID: p.SystemID,
	}
}

func (p *BindTransmitter) Resp() Packet {
	return &BindTransmitterResp{
		Header:   p.makeHeader(),
		SystemID: p.SystemID,
	}
}

func (p *BroadcastSM) Resp() Packet {
	return &BroadcastSMResp{
		Header:    p.makeHeader(),
		MessageID: p.MessageID,
	}
}

func (p *CancelBroadcastSM) Resp() Packet {
	return &CancelBroadcastSMResp{
		Header: p.makeHeader(),
	}
}

func (p *CancelSM) Resp() Packet {
	return &CancelSMResp{
		Header: p.makeHeader(),
	}
}

func (p *DataSM) Resp() Packet {
	return &DataSMResp{
		Header: p.makeHeader(),
	}
}

func (p *DeliverSM) Resp() Packet {
	return &DeliverSMResp{
		Header: p.makeHeader(),
	}
}

func (p *EnquireLink) Resp() Packet {
	return &EnquireLinkResp{
		Header: p.makeHeader(),
	}
}

func (p *QueryBroadcastSM) Resp() Packet {
	return &QueryBroadcastSMResp{
		Header:    p.makeHeader(),
		MessageID: p.MessageID,
	}
}

func (p *QuerySM) Resp() Packet {
	return &QuerySMResp{
		Header: p.makeHeader(),
	}
}

func (p *ReplaceSM) Resp() Packet {
	return &ReplaceSMResp{
		Header: p.makeHeader(),
	}
}

func (p *SubmitMulti) Resp() Packet {
	return &SubmitMultiResp{
		Header: p.makeHeader(),
	}
}

func (p *SubmitSM) Resp() Packet {
	return &SubmitSMResp{
		Header: p.makeHeader(),
	}
}

func (p *Unbind) Resp() Packet {
	return &UnbindResp{
		Header: p.makeHeader(),
	}
}
