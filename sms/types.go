package sms

type directionSetter interface {
	setDirection(Direction)
}

type SMPPMarshaller interface {
	MarshalSMPP() (interface{}, error)
}

type SMPPUnmarshaler interface {
	UnmarshalSMPP(interface{}) error
}
