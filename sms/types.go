package sms

type SMPPMarshaler interface {
	MarshalSMPP() (interface{}, error)
}

type SMPPUnmarshaler interface {
	UnmarshalSMPP(interface{}) error
}
