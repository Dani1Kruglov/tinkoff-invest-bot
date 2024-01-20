package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"log"
	"sync"
	"tinkoff-investment-bot/internal/bot/model"
	ms "tinkoff-investment-bot/internal/model/settings"
	t "tinkoff-investment-bot/internal/model/tracker"
	"tinkoff-investment-bot/internal/services/instruments/shares"
	"tinkoff-investment-bot/internal/tracking"
)

func InlineKeyBoardHandler(tinkoffInvestBot *model.Bot, settings *ms.Settings, client *investgo.Client, cacheCommand *sync.Map) {
	if len(tinkoffInvestBot.Update.CallbackQuery.Data) > 7 && tinkoffInvestBot.Update.CallbackQuery.Data[:7] == "command" {
		inlineKeyBoardMainMenu(tinkoffInvestBot, settings, client, cacheCommand)

	} else if lastCommandByChatIDAndTicker, ok := cacheCommand.LoadAndDelete(tinkoffInvestBot.Update.FromChat().ID); ok {
		if lastCommandByChatID, ok := lastCommandByChatIDAndTicker.(string); ok {
			switch lastCommandByChatID[4:] {
			case "3SaveOrNot":
				analysisCallbackDataOnYesOrNo(tinkoffInvestBot, settings, cacheCommand, lastCommandByChatID[:4]+"3WPrice",
					"Отслеживать изменение цены на 10% от текущей или укажите свою начальную цену: ", &tgbotapi.InlineKeyboardMarkup{
						InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
							{
								tgbotapi.NewInlineKeyboardButtonData("Да", "yes"),
							},
						},
					})
			case "3WPrice":
				responses := shares.AddShareToListOfTracked(t.NewTracker(client), settings, lastCommandByChatID[:4], tinkoffInvestBot.Update.FromChat().ID, "")
				printResponseToBot(tinkoffInvestBot, responses,
					settings.Logger, nil)
				mainMenu(tinkoffInvestBot)

			}
		}
	} else {
		inlineKeyBoardStart(tinkoffInvestBot, client)
		cacheCommand.Store(tinkoffInvestBot.Update.FromChat().ID, "token")
	}
}

func analysisCallbackDataOnYesOrNo(tinkoffInvestBot *model.Bot, settings *ms.Settings, cacheCommand *sync.Map, cacheValue string, msgForYes string, numericalKeyboard *tgbotapi.InlineKeyboardMarkup) {
	if tinkoffInvestBot.Update.CallbackQuery.Data == "yes" {
		cacheCommand.Store(tinkoffInvestBot.Update.FromChat().ID, cacheValue)
		printResponseToBot(tinkoffInvestBot, []string{msgForYes},
			settings.Logger, numericalKeyboard)
	} else {
		mainMenu(tinkoffInvestBot)
	}
}

func inlineKeyBoardStart(tinkoffInvestBot *model.Bot, client *investgo.Client) {
	switch tinkoffInvestBot.Update.CallbackQuery.Data {
	case "1":
		if client != nil {
			mainMenu(tinkoffInvestBot)
		} else {
			msg := tgbotapi.NewMessage(tinkoffInvestBot.Update.FromChat().ID, "")
			msg.Text = "Введите 'токен для чтения' своего брокерского счета:"
			if _, err := tinkoffInvestBot.Api.Send(msg); err != nil {
				log.Fatalln(err.Error())
			}
		}
		break
	case "2":
		break
	}
}

func inlineKeyBoardMainMenu(tinkoffInvestBot *model.Bot, settings *ms.Settings, client *investgo.Client, cacheCommand *sync.Map) {
	if client != nil {
		responses := tracking.TrackByTinkoffToken(settings, t.NewTracker(client), tinkoffInvestBot.Update.FromChat().ID, tinkoffInvestBot.Update.CallbackQuery.Data[8:])
		printResponseToBot(tinkoffInvestBot, responses, settings.Logger, nil)
		if tinkoffInvestBot.Update.CallbackQuery.Data[8:] != "2" && tinkoffInvestBot.Update.CallbackQuery.Data[8:] != "3" {
			mainMenu(tinkoffInvestBot)
			return
		}
		cacheCommand.Store(tinkoffInvestBot.Update.FromChat().ID, tinkoffInvestBot.Update.CallbackQuery.Data[8:])
	}
}
