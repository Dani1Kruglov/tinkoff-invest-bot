package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tinkoff-investment-bot/internal/bot/commands"
	"tinkoff-investment-bot/internal/bot/model"
)

func CommandHandler(tinkoffInvestBot *model.Bot) {
	msg := tgbotapi.NewMessage(tinkoffInvestBot.Update.FromChat().ID, "")
	switch tinkoffInvestBot.Update.Message.Command() {
	case "start":
		msg.Text, msg.ReplyMarkup = commands.ViewCommandStart()
		break
	default:
		msg.Text = commands.ViewCommandDefault()
		break
	}

	if _, err := tinkoffInvestBot.Api.Send(msg); err != nil {
		log.Fatalln(err.Error())
	}
}
