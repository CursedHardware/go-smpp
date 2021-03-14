package pdu

func ReadSequence(packet Packet) int32 {
	if packet == nil {
		return -1
	}
	return packet.getSequence()
}

func WriteSequence(packet Packet, sequence int32) {
	packet.setSequence(sequence)
}

func ReadCommandStatus(packet Packet) CommandStatus {
	if packet == nil {
		return 0
	}
	return packet.getCommandStatus()
}
