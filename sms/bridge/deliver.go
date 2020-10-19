package bridge

import (
	"bytes"

	"github.com/VoiceGateway/go-smpp/pdu"
	"github.com/VoiceGateway/go-smpp/sms"
)

func ToDeliverSM(deliver *sms.Deliver) (sm *pdu.DeliverSM, err error) {
	var message pdu.ShortMessage
	if deliver.Flags.UDHIndicator {
		message.UDHeader = pdu.UserDataHeader{}
	}
	_, err = message.ReadFrom(bytes.NewReader(deliver.UserData))
	sm = &pdu.DeliverSM{
		SourceAddr: pdu.Address{
			TON: deliver.OriginatingAddress.TON,
			NPI: deliver.OriginatingAddress.NPI,
			No:  deliver.OriginatingAddress.No,
		},
		ESMClass: pdu.ESMClass{
			MessageType:  deliver.Flags.MessageType.Type(),
			UDHIndicator: deliver.Flags.UDHIndicator,
			ReplyPath:    deliver.Flags.ReplyPath,
		},
		ProtocolID: deliver.ProtocolIdentifier,
		Message:    message,
	}
	return
}
