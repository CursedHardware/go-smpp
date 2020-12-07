package bridge

import (
	"bytes"

	"github.com/M2MGateway/go-smpp/pdu"
	"github.com/M2MGateway/go-smpp/sms"
)

func ToSubmit(sm *pdu.SubmitSM) (submit *sms.Submit, err error) {
	var userData bytes.Buffer
	_, _ = sm.Message.UDHeader.WriteTo(&userData)
	userData.Write(sm.Message.Message)
	submit = &sms.Submit{
		Flags: sms.SubmitFlags{
			ReplyPath:               sm.ESMClass.ReplyPath,
			UserDataHeaderIndicator: sm.Message.UDHeader != nil,
		},
		DestinationAddress: sms.Address{
			NPI: sm.DestAddr.NPI,
			TON: sm.DestAddr.TON,
			No:  sm.DestAddr.No,
		},
		ProtocolIdentifier: sm.ProtocolID,
		DataCoding:         byte(sm.Message.DataCoding),
		UserData:           userData.Bytes(),
	}
	return
}
