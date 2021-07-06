package pdu

import "fmt"

// RegisteredDelivery see SMPP v5, section 4.7.21 (130p)
type RegisteredDelivery struct {
	MCDeliveryReceipt           byte // ___ _ __ **
	SMEOriginatedAcknowledgment byte // ___ _ ** __
	IntermediateNotification    bool // ___ * __ __
	Reserved                    byte // *** _ __ __
}

func (r RegisteredDelivery) ReadByte() (c byte, err error) {
	c |= r.MCDeliveryReceipt & 0b11
	c |= r.SMEOriginatedAcknowledgment & 0b11 << 2
	c |= getBool(r.IntermediateNotification) << 4
	c |= r.Reserved & 0b111 << 5
	return
}

func (r *RegisteredDelivery) WriteByte(c byte) error {
	r.MCDeliveryReceipt = c & 0b11
	r.SMEOriginatedAcknowledgment = c >> 2 & 0b11
	r.IntermediateNotification = c>>4&0b1 == 1
	r.Reserved = c >> 5 & 0b111
	return nil
}

func (r RegisteredDelivery) String() string {
	c, _ := r.ReadByte()
	return fmt.Sprintf("%08b", c)
}
