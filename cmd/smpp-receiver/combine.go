package main

import (
	"log"
	"time"

	"github.com/M2MGateway/go-smpp/pdu"
)

func makeCombineMultipartDeliverSM(device *Device, hook func(*Payload)) func(*pdu.DeliverSM) {
	return pdu.CombineMultipartDeliverSM(func(delivers []*pdu.DeliverSM) {
		var mergedMessage string
		for _, sm := range delivers {
			if sm.Message.DataCoding == 0x00 && device.Workaround == "SMG4000" {
				mergedMessage += string(sm.Message.Message)
			} else if message, err := sm.Message.Parse(); err == nil {
				mergedMessage += message
			}
		}
		source := delivers[0].SourceAddr
		target := delivers[0].DestAddr
		log.Println(device, source, "->", target)
		go hook(&Payload{
			SMSC:        device.SMSC,
			SystemID:    device.SystemID,
			SystemType:  device.SystemType,
			Owner:       device.Owner,
			Phone:       device.Phone,
			Extra:       device.Extra,
			Source:      source.String(),
			Target:      target.String(),
			Message:     mergedMessage,
			DeliverTime: time.Now(),
		})
	})
}
