package sms

type directionSetter interface {
	setDirection(Direction)
}

type SMPPMarshaler interface {
	MarshalSMPP() (interface{}, error)
}

type SMPPUnmarshaler interface {
	UnmarshalSMPP(interface{}) error
}
