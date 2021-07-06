package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

const (
	sender   = "sender"
	receiver = "receiver"
	submit   = "submit"
)

func makeKeyboard() *tgbotapi.InlineKeyboardMarkup {
	onTargetRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Sender", sender),
		tgbotapi.NewInlineKeyboardButtonData("Receiver", receiver),
	)
	onSubmitRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Submit", submit),
	)
	markup := tgbotapi.NewInlineKeyboardMarkup(onTargetRow, onSubmitRow)
	return &markup
}
