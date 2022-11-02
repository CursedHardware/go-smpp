package sms

type directionSetter interface {
	setDirection(Direction)
}

type SMPPMarshaller interface {
	MarshalSMPP() (any, error)
}

type SMPPUnmarshaler interface {
	UnmarshalSMPP(any) error
}
