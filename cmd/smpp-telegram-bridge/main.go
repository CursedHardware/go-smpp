package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	. "github.com/go-telegram-bot-api/telegram-bot-api"
)

var configure Configuration

func init() {
	var confPath string
	flag.StringVar(&confPath, "conf", "configure.json", "configure file-path")
	if data, err := ioutil.ReadFile(confPath); err != nil {
		log.Panic(err)
	} else if err = json.Unmarshal(data, &configure); err != nil {
		log.Panic(err)
	}
}

func main() {
	bot, err := NewBotAPIWithClient(configure.Token, http.DefaultClient)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %q", bot.Self.UserName)
	updates, err := bot.GetUpdatesChan(UpdateConfig{Timeout: 60})
	if err != nil {
		log.Panic(err)
	}
	for update := range updates {
		if update.Message != nil {
			go handleUpdate(bot, update)
		}
	}
}

func handleUpdate(bot *BotAPI, update Update) {
	message := update.Message
	chat := message.Chat
	if !chat.IsPrivate() {
		_, _ = bot.LeaveChat(chat.ChatConfig())
		return
	} else if !configure.isAllowedUserID(chat.ID) {
		msg := NewMessage(chat.ID, fmt.Sprintf("Your User ID is %d, Please contact the owner", chat.ID))
		msg.ReplyToMessageID = message.MessageID
		_, _ = bot.Send(msg)
		return
	}
}
