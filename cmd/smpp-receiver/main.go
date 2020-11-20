package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"sync"
	"time"

	"github.com/VoiceGateway/go-smpp"
	"github.com/VoiceGateway/go-smpp/pdu"
)

var configure Configuration
var mutex sync.Mutex

func init() {
	configure.HookMode = "event"
	var confPath string
	flag.StringVar(&confPath, "conf", "configure.json", "configure file-path")
	if data, err := ioutil.ReadFile(confPath); err != nil {
		log.Fatal(err)
	} else if err = json.Unmarshal(data, &configure); err != nil {
		log.Fatal(err)
	}
}

func main() {
	hook := runProgram
	if configure.HookMode == "ndjson" {
		hook = runProgramWithEvent()
	}
	for _, device := range configure.Devices {
		fillAccount(&device)
		go func(device Account) {
			for {
				connect(device, hook)
			}
		}(device)
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
	if account.Extra == nil {
		account.Extra = configure.DefaultAccount.Extra
	}
}

//goland:noinspection GoUnhandledErrorResult
func connect(device Account, hook func(*Payload)) {
	parent, err := net.Dial("tcp", device.SMSC)
	if err != nil {
		log.Fatalln(err)
	}
	ctx := context.Background()
	conn := smpp.NewConn(ctx, parent)
	conn.ReadTimeout = time.Second * 30
	go conn.Watch()
	defer conn.Close()
	resp, err := conn.Submit(ctx, &pdu.BindTransceiver{
		SystemID:   device.SystemID,
		Password:   device.Password,
		SystemType: device.SystemType,
		Version:    pdu.SMPPVersion34,
	})
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
			Message:     merged,
			DeliverTime: time.Now(),
			Extra:       device.Extra,
		})
	})
	for {
		packet := <-conn.PDU()
		switch p := packet.(type) {
		case *pdu.DeliverSM:
			addDeliverSM(p)
			_ = conn.Send(p.Resp())
		case pdu.Responsable:
			_ = conn.Send(p.Resp())
		}
	}
}

//goland:noinspection GoUnhandledErrorResult
func runProgram(message *Payload) {
	mutex.Lock()
	defer mutex.Unlock()
	log.Printf("%s @ %s | %s -> %s", message.SMSC, message.SystemID, message.Source, message.Target)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*15)
	defer cancel()
	cmd := exec.CommandContext(ctx, configure.Hook)
	go func(stdin io.WriteCloser, err error) {
		if err != nil {
			log.Fatal(err)
		}
		defer stdin.Close()
		_ = json.NewEncoder(stdin).Encode(message)
	}(cmd.StdinPipe())
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatal(message, err, "\n", string(output))
	}
}

func runProgramWithEvent() func(*Payload) {
	cmd := exec.Command(configure.Hook)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	go func() {
		if output, err := cmd.CombinedOutput(); err != nil {
			log.Fatal(err, "\n", string(output))
		}
	}()
	return func(message *Payload) {
		log.Printf("%s @ %s | %s -> %s", message.SMSC, message.SystemID, message.Source, message.Target)
		_ = json.NewEncoder(stdin).Encode(message)
		_, _ = fmt.Fprintln(stdin)
	}
}
