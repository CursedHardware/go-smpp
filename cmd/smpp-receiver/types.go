package main

import (
	"encoding/json"
	"time"
)

type Configuration struct {
	Hook           string    `json:"hook"`
	HookMode       string    `json:"hook_mode"`
	DefaultAccount Account   `json:"default_account"`
	Devices        []Account `json:"devices"`
}

//goland:noinspection ALL
type Account struct {
	SMSC       string          `json:"smsc"`
	SystemID   string          `json:"system_id"`
	Password   string          `json:"password"`
	SystemType string          `json:"system_type"`
	Extra      json.RawMessage `json:"extra"`
}

//goland:noinspection ALL
type Payload struct {
	SMSC        string          `json:"smsc"`
	SystemID    string          `json:"system_id"`
	SystemType  string          `json:"system_type"`
	Source      string          `json:"source"`
	Target      string          `json:"target"`
	Message     string          `json:"message"`
	DeliverTime time.Time       `json:"deliver_time"`
	Extra       json.RawMessage `json:"extra,omitempty"`
}
