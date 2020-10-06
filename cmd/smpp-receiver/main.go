package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/NiceLabs/go-smpp"
	"github.com/NiceLabs/go-smpp/pdu"
)

var configure = new(Configuration)

func init() {
	configFile, err := ioutil.ReadFile("configure.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(configFile, configure)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	for _, device := range configure.Devices {
		fillAccount(&device)
		go connect(device, runProgram)
	}
	select {}
}

func fillAccount(account *Account) {
	if account.SMSC == "" {
		account.SMSC = configure.DefaultAccount.SMSC
	}
	if account.SystemID == "" {
		account.SystemID = configure.DefaultAccount.SystemID
	}
	if account.SystemType == "" {
		account.SystemType = configure.DefaultAccount.SystemType
	}
	if account.Password == "" {
		account.Password = configure.DefaultAccount.Password
	}
	if account.BindType == "" {
		account.BindType = configure.DefaultAccount.BindType
	}
}

func connect(device Account, hook func(*Payload)) {
	parent, err := net.Dial("tcp", device.SMSC)
	if err != nil {
		log.Fatalln(err)
	}
	ctx := context.Background()
	conn := smpp.NewConn(ctx, parent)
	conn.ReadTimeout = time.Second * 30
	go conn.Watch()
	var request pdu.Responsable
	switch device.BindType {
	case "", "receiver":
		request = &pdu.BindReceiver{
			SystemID:   device.SystemID,
			Password:   device.Password,
			SystemType: device.SystemType,
			Version:    pdu.SMPPVersion50,
		}
	case "transceiver":
		request = &pdu.BindTransceiver{
			SystemID:   device.SystemID,
			Password:   device.Password,
			SystemType: device.SystemType,
			Version:    pdu.SMPPVersion50,
		}
	case "transmitter":
		request = &pdu.BindTransmitter{
			SystemID:   device.SystemID,
			Password:   device.Password,
			SystemType: device.SystemType,
			Version:    pdu.SMPPVersion50,
		}
	default:
		log.Fatalln("unsupported bind type")
	}
	resp, err := conn.Submit(ctx, request)
	if err != nil {
		log.Fatalln(err)
	} else if status := pdu.ReadCommandStatus(resp); status != 0 {
		log.Fatalln(status)
	}
	log.Printf("Connected %s @ %s", device.SMSC, device.SystemID)
	go conn.EnquireLink(time.Second*5, time.Minute)
	addDeliverSM := pdu.CombineMultipartDeliverSM(func(pdu []*pdu.DeliverSM) {
		var merged string
		for _, sm := range pdu {
			message, err := sm.Message.Parse()
			if err != nil {
				continue
			}
			merged += message
		}
		hook(&Payload{
			SMSC:        device.SMSC,
			SystemID:    device.SystemID,
			SystemType:  device.SystemType,
			Source:      pdu[0].SourceAddr.No,
			Target:      pdu[0].DestAddr.No,
			Message:     strings.ReplaceAll(merged, "\x7f\x7f ", "\n"),
			DeliverTime: time.Now(),
		})
	})
	for {
		packet := <-conn.PDU()
		switch p := packet.(type) {
		case *pdu.DeliverSM:
			addDeliverSM(p)
			_ = conn.Send(ctx, p.Resp())
		case pdu.Responsable:
			_ = conn.Send(ctx, p.Resp())
		}
	}
}

//goland:noinspection GoUnhandledErrorResult
func runProgram(message *Payload) {
	log.Printf("%s @ %s | %s -> %s", message.SMSC, message.SystemID, message.Source, message.Target)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*15)
	defer cancel()
	cmd := exec.CommandContext(ctx, configure.Hook)
	go func(stdin io.WriteCloser, err error) {
		if err != nil {
			log.Fatal(err)
		}
		defer stdin.Close()
		err = json.NewEncoder(stdin).Encode(message)
		if err != nil {
			log.Fatal(err)
		}
	}(cmd.StdinPipe())
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatal(message, err, "\n", string(output))
	}
}
