package pdu

import "fmt"

// RegisteredDelivery see SMPP v5, section 4.7.21 (130p)
type RegisteredDelivery struct {
	MCDeliveryReceipt           byte
	SMEOriginatedAcknowledgment byte
	IntermediateNotification    bool
}

func (r RegisteredDelivery) ReadByte() (c byte, err error) {
	c |= r.MCDeliveryReceipt & 0b11
	c |= r.SMEOriginatedAcknowledgment & 0b11 << 2
	if r.IntermediateNotification {
		c |= 1 << 4
	}
	return
}

func (r *RegisteredDelivery) WriteByte(c byte) error {
	r.MCDeliveryReceipt = c & 0b11
	r.SMEOriginatedAcknowledgment = c >> 2 & 0b11
	r.IntermediateNotification = c>>4&0b1 == 1
	return nil
}

func (r RegisteredDelivery) String() string {
	c, _ := r.ReadByte()
	return fmt.Sprintf("%08b", c)
}
