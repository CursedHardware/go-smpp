package main

import (
	"time"
)

type Configuration struct {
	Hook           string    `json:"hook"`
	DefaultAccount Account   `json:"default_account,omitempty"`
	Devices        []Account `json:"devices"`
}

//goland:noinspection ALL
type Account struct {
	SMSC       string `json:"smsc"`
	SystemID   string `json:"system_id"`
	Password   string `json:"password"`
	SystemType string `json:"system_type"`
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
