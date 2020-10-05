package main

import (
	"time"
)

type Configuration struct {
	DefaultAccount Account `json:",omitempty"`
	Devices        []Account
	Hook           string
}

//goland:noinspection ALL
type Account struct {
	SMSC       string
	SystemID   string
	Password   string
	SystemType string
	BindType   string `json:"bind_type"`
}

//goland:noinspection ALL
type Payload struct {
	SMSC        string    `json:"smsc"`
	SystemID    string    `json:"system_id"`
	SystemType  string    `json:"system_type"`
	Source      string    `json:"source"`
	Target      string    `json:"target"`
	Message     string    `json:"message"`
	DeliverTime time.Time `json:"deliver_time"`
}
