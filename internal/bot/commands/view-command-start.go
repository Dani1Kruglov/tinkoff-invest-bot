package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("1 - только для отслеживания информации по ценным бумагам ", "1"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("2 - для выполнения операций", "2"),
	),
)

func ViewCommandStart() (string, *tgbotapi.InlineKeyboardMarkup) {
	return "Хотите использовать бота для: ", &numericKeyboard
}
