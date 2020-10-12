package sms

type MessageType byte

const (
	MessageTypeDeliver MessageType = iota
	MessageTypeDeliverReport
	MessageTypeSubmitReport
	MessageTypeSubmit
	MessageTypeStatusReport
	MessageTypeCommand
)

type Direction int

const (
	MT Direction = iota
	MO
)

func (t *MessageType) Set(kind byte, dir Direction) {
	*t = MessageType(kind<<1 | byte(dir))
}

func (t MessageType) Type() MessageType {
	return t >> 1
}

func (t MessageType) Direction() Direction {
	return Direction(t & 0x01)
}

func (t MessageType) String() string {
	switch t {
	case MessageTypeDeliverReport:
		return "SMS-DELIVER-REPORT"
	case MessageTypeDeliver:
		return "SMS-DELIVER"
	case MessageTypeSubmit:
		return "SMS-SUBMIT"
	case MessageTypeSubmitReport:
		return "SMS-SUBMIT-REPORT"
	case MessageTypeCommand:
		return "SMS-COMMAND"
	case MessageTypeStatusReport:
		return "SMS-STATUS-REPORT"
	}
	return "[UNKNOWN]"
}
