package pdu

import (
	"encoding/binary"
	"io"
)

// CommandID see SMPP v5, section 4.7.5 (115p)
type CommandID uint32

// CommandStatus see SMPP v5, section 4.7.6 (116p)
type CommandStatus uint32

type Header struct {
	CommandLength uint32
	CommandID     CommandID
	CommandStatus CommandStatus
	Sequence      int32
}

func (h *Header) getCommandStatus() CommandStatus {
	return h.CommandStatus
}

func (h *Header) getSequence() int32 {
	return h.Sequence
}

func (h *Header) setSequence(sequence int32) {
	h.Sequence = sequence
}

func (h *Header) makeHeader() Header {
	return Header{Sequence: h.Sequence}
}

func readHeaderFrom(r io.Reader, header *Header) (err error) {
	err = binary.Read(r, binary.BigEndian, header)
	if err == nil && (header.CommandLength < 16 || header.CommandLength > 0x10000) {
		err = ErrInvalidCommandLength
	}
	return
}
