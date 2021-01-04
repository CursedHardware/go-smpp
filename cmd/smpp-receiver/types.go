package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/M2MGateway/go-smpp/pdu"
)

type Configuration struct {
	Hook     string    `json:"hook"`
	HookMode string    `json:"hook_mode"`
	Devices  []*Device `json:"devices"`
}

//goland:noinspection ALL
type Device struct {
	SMSC             string               `json:"smsc"`
	SystemID         string               `json:"system_id"`
	Password         string               `json:"password"`
	SystemType       string               `json:"system_type"`
	Version          pdu.InterfaceVersion `json:"version"`
	BindMode         string               `json:"bind_mode"`
	Owner            string               `json:"owner"`
	Phone            string               `json:"phone"`
	Extra            json.RawMessage      `json:"extra"`
	Workaround       string               `json:"workaround"`
	KeepAliveTick    time.Duration        `json:"keepalive_tick"`
	KeepAliveTimeout time.Duration        `json:"keepalive_timeout"`
}

func (d *Device) String() string {
	return fmt.Sprintf("%s @ %s", d.SMSC, d.SystemID)
}

func (d *Device) Binder() pdu.Responsable {
	switch d.BindMode {
	case "receiver":
		return &pdu.BindReceiver{
			SystemID:   d.SystemID,
			Password:   d.Password,
			SystemType: d.SystemType,
			Version:    d.Version,
		}
	case "transceiver":
		return &pdu.BindTransceiver{
			SystemID:   d.SystemID,
			Password:   d.Password,
			SystemType: d.SystemType,
			Version:    d.Version,
		}
	}
	return nil
}

//goland:noinspection ALL
type Payload struct {
	SMSC        string          `json:"smsc"`
	SystemID    string          `json:"system_id"`
	SystemType  string          `json:"system_type"`
	Source      string          `json:"source,omitempty"`
	Target      string          `json:"target,omitempty"`
	Message     string          `json:"message"`
	DeliverTime time.Time       `json:"deliver_time"`
	Owner       string          `json:"owner,omitempty"`
	Phone       string          `json:"phone,omitempty"`
	Extra       json.RawMessage `json:"extra,omitempty"`
}
