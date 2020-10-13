package sms

type MessageType byte

const (
	MessageTypeDeliver       MessageType = iota // 00 0 MT
	MessageTypeDeliverReport                    // 00 1 MO
	MessageTypeSubmitReport                     // 01 0 MT
	MessageTypeSubmit                           // 01 1 MO
	MessageTypeStatusReport                     // 10 0 MT
	MessageTypeCommand                          // 10 1 MO
)

type Direction int

const (
	MT Direction = iota
	MO
)

func (t *MessageType) Set(kind byte, dir Direction) {
	*t = MessageType(kind<<1 | byte(dir))
}

func (t MessageType) Type() byte {
	return byte(t >> 1)
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
