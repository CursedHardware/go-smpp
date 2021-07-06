package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/M2MGateway/go-smpp/cmd/smpp-telegram-bridge/templates"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
	db, err := gorm.Open(sqlite.Open("bridge.db?cache=shared&_busy_timeout=5000"), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	bot, err := tgbotapi.NewBotAPIWithClient(configure.Token, http.DefaultClient)
	if err != nil {
		log.Panic(err)
	}
	err = db.AutoMigrate(new(ShortMessage))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %q", bot.Self.UserName)
	updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})
	if err != nil {
		log.Panic(err)
	}
	for update := range updates {
		switch {
		case update.Message != nil:
			go handleMessage(db, bot, update.Message)
		case update.EditedMessage != nil:
			go handleEditedMessage(db, bot, update.EditedMessage)
		case update.CallbackQuery != nil:
			go handleCallbackQuery(db, bot, update.CallbackQuery)
		}
	}
}

func handleMessage(db *gorm.DB, bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	switch {
	case !message.Chat.IsPrivate():
		_, _ = bot.LeaveChat(message.Chat.ChatConfig())
	case !configure.isAllowedUserID(message.From.ID):
		_, _ = bot.Send(tgbotapi.MessageConfig{
			BaseChat:  tgbotapi.BaseChat{ChatID: message.Chat.ID, ReplyToMessageID: message.MessageID},
			ParseMode: "HTML",
			Text:      templates.Must(templates.Unauthorized(message.From)),
		})
	case message.Text != "":
		slave, _ := bot.Send(tgbotapi.MessageConfig{
			BaseChat:  tgbotapi.BaseChat{ChatID: message.Chat.ID, ReplyToMessageID: message.MessageID, ReplyMarkup: makeKeyboard()},
			ParseMode: "HTML",
			Text:      templates.Must(templates.Reply(&templates.ReplyData{Segments: getSegments(message.Text)})),
		})
		row := &ShortMessage{
			ChatID:          message.Chat.ID,
			ParentMessageID: message.MessageID,
			SlaveMessageID:  slave.MessageID,
		}
		db.Create(row).Commit()
	}
}

func handleEditedMessage(db *gorm.DB, bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	store := &ShortMessage{ChatID: message.Chat.ID, ParentMessageID: message.MessageID}
	if err := db.Model(store).Take(store).Error; err != nil {
		log.Println(err)
		return
	}
	_, _ = bot.Send(tgbotapi.EditMessageTextConfig{
		BaseEdit:  tgbotapi.BaseEdit{ChatID: message.Chat.ID, ReplyMarkup: makeKeyboard()},
		ParseMode: "HTML",
		Text:      templates.Must(templates.Reply(&templates.ReplyData{Segments: getSegments(message.Text)})),
	})
}

func handleCallbackQuery(db *gorm.DB, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
	msg := &ShortMessage{ChatID: query.Message.Chat.ID, ParentMessageID: query.Message.MessageID}
	if err := db.Model(msg).Take(msg).Error; err != nil {
		log.Println(err)
		return
	}
	switch query.Data {
	case sender:
	case receiver:
	case submit:
		if !msg.Ready() {
			_, _ = bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
				CallbackQueryID: query.ID,
				Text:            "Not Ready",
				ShowAlert:       true,
			})
		} else {
			_, _ = bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
				CallbackQueryID: query.ID,
				Text:            "Ready",
				ShowAlert:       true,
			})
		}
	}
}
