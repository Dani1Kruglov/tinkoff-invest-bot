package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"log"
	"tinkoff-investment-bot/internal/bot/model"
	ms "tinkoff-investment-bot/internal/model/settings"
	"tinkoff-investment-bot/internal/tracking"
)

func InlineKeyBoardHandler(tinkoffInvestBot *model.Bot, settings *ms.Settings, client *investgo.Client) {
	if len(tinkoffInvestBot.Update.CallbackQuery.Data) > 7 && tinkoffInvestBot.Update.CallbackQuery.Data[:7] == "command" {
		inlineKeyBoardMainMenu(tinkoffInvestBot, settings, client)
	} else {
		inlineKeyBoardStart(tinkoffInvestBot, client)
	}
}

func inlineKeyBoardStart(tinkoffInvestBot *model.Bot, client *investgo.Client) {
	switch tinkoffInvestBot.Update.CallbackQuery.Data {
	case "1":
		if client != nil {
			MainMenu(tinkoffInvestBot)
		} else {
			msg := tgbotapi.NewMessage(tinkoffInvestBot.Update.FromChat().ID, "")
			msg.Text = "Введите 'токен для чтения' своего брокерского счета (token=<Your token>):"
			if _, err := tinkoffInvestBot.Api.Send(msg); err != nil {
				log.Fatalln(err.Error())
			}
		}
		break
	case "2":
		break
	}
}

func inlineKeyBoardMainMenu(tinkoffInvestBot *model.Bot, settings *ms.Settings, client *investgo.Client) {
	if client != nil {
		responses := tracking.TrackByTinkoffToken(settings, client, tinkoffInvestBot.Update.FromChat().ID, tinkoffInvestBot.Update.CallbackQuery.Data[8:])
		printRespToBot(tinkoffInvestBot, responses)
		if tinkoffInvestBot.Update.CallbackQuery.Data[8:] != "2" {
			MainMenu(tinkoffInvestBot)
		}
	}
}
