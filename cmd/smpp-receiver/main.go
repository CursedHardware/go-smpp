package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/M2MGateway/go-smpp"
	"github.com/M2MGateway/go-smpp/pdu"
	"github.com/imdario/mergo"
	. "github.com/xeipuuv/gojsonschema"
)

var configure Configuration
var mutex sync.Mutex

//go:embed schema.json
var schemaFile []byte

func init() {
	var confPath string
	flag.StringVar(&confPath, "c", "configure.json", "configure file-path")

	if data, err := ioutil.ReadFile(confPath); err != nil {
		log.Fatalln(err)
	} else if result, _ := Validate(NewBytesLoader(schemaFile), NewBytesLoader(data)); !result.Valid() {
		for _, desc := range result.Errors() {
			log.Println(desc)
		}
		log.Fatalln("invalid configuration")
	} else {
		_ = json.Unmarshal(data, &configure)
	}
	_ = mergo.Merge(configure.Devices[0], &Device{
		Version:          pdu.SMPPVersion34,
		BindMode:         "receiver",
		KeepAliveTick:    time.Millisecond * 500,
		KeepAliveTimeout: time.Second,
	})
	for i := 1; i < len(configure.Devices); i++ {
		_ = mergo.Merge(configure.Devices[i], configure.Devices[i-1])
	}
}

func main() {
	hook := runProgramWithEvent
	if configure.HookMode == "ndjson" {
		hook = runProgramWithStream()
	}
	for _, device := range configure.Devices {
		go connect(device, hook)
	}
	select {}
}

//goland:noinspection GoUnhandledErrorResult
func connect(device *Device, hook func(*Payload)) {
	conn, err := smpp.DialTimeout(device.SMSC, time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	conn.ReadTimeout = time.Second
	conn.WriteTimeout = time.Second
	defer conn.Close()
	if resp, err := conn.Submit(context.Background(), device.Binder()); err != nil {
		log.Fatalln(device, err)
	} else if status := pdu.ReadCommandStatus(resp); status != 0 {
		log.Fatalln(device, status)
	} else {
		log.Println(device, "Connected")
		go conn.EnquireLink(device.KeepAliveTick, device.KeepAliveTimeout)
	}
	addDeliverSM := makeCombineMultipartDeliverSM(device, hook)
	for packet := range conn.PDU() {
		switch p := packet.(type) {
		case *pdu.DeliverSM:
			addDeliverSM(p)
			_ = conn.Send(p.Resp())
		case pdu.Responsable:
			_ = conn.Send(p.Resp())
		}
	}
	log.Println(device, "Disconnected")
	time.Sleep(time.Second)
	go connect(device, hook)
}
