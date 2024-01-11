package model

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	Api    *tgbotapi.BotAPI
	Update tgbotapi.Update
}

func New(api *tgbotapi.BotAPI) *Bot {
	return &Bot{
		Api: api,
	}
}
