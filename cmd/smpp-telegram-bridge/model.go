package main

import "gorm.io/gorm"

type ShortMessage struct {
	gorm.Model
	ChatID          int64 `gorm:"index"`
	SlaveMessageID  int   `gorm:"index"`
	ParentMessageID int   `gorm:"index"`
	Sender          string
	Receiver        string
}

func (m *ShortMessage) Ready() bool {
	return m.Sender != "" && m.Receiver != ""
}
