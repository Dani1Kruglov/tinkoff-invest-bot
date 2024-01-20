package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"go.uber.org/zap"
	"log"
	"strings"
	"sync"
	"tinkoff-investment-bot/internal/bot/model"
	"tinkoff-investment-bot/internal/connect/tinkoff-client"
	ms "tinkoff-investment-bot/internal/model/settings"
	t "tinkoff-investment-bot/internal/model/tracker"
	"tinkoff-investment-bot/internal/services/instruments/shares"
	"tinkoff-investment-bot/internal/tracking"
)

func MessageHandler(tinkoffInvestBot *model.Bot, settings *ms.Settings, client *investgo.Client, cancel *context.CancelFunc, cacheCommand *sync.Map) {
	if lastCommandByChatIDAndTicker, ok := cacheCommand.LoadAndDelete(tinkoffInvestBot.Update.FromChat().ID); ok {
		if lastCommandByChatID, ok := lastCommandByChatIDAndTicker.(string); ok {
			switch {
			case lastCommandByChatID == "2" && client != nil:
				responses := tracking.GetShare(settings, t.NewTracker(client), strings.ToUpper(tinkoffInvestBot.Update.Message.Text))
				printResponseToBot(tinkoffInvestBot, responses, settings.Logger, nil)
				mainMenu(tinkoffInvestBot)
			case lastCommandByChatID == "3" && client != nil:
				responses := shares.GetShareForFavoriteList(t.NewTracker(client), settings, strings.ToUpper(tinkoffInvestBot.Update.Message.Text))
				printResponseToBot(tinkoffInvestBot, responses, settings.Logger, &tgbotapi.InlineKeyboardMarkup{
					InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
						{
							tgbotapi.NewInlineKeyboardButtonData("Да", "yes"),
						},
						{
							tgbotapi.NewInlineKeyboardButtonData("Нет", "no"),
						},
					},
				})
				cacheCommand.Store(tinkoffInvestBot.Update.FromChat().ID, strings.ToUpper(tinkoffInvestBot.Update.Message.Text)+"3SaveOrNot")
			case lastCommandByChatID[4:] == "3WPrice":
				responses := shares.AddShareToListOfTracked(t.NewTracker(client), settings, lastCommandByChatID[:4], tinkoffInvestBot.Update.FromChat().ID, tinkoffInvestBot.Update.Message.Text)
				printResponseToBot(tinkoffInvestBot, responses,
					settings.Logger, nil)

			case lastCommandByChatID == "token":
				client, *cancel = tinkoff_client.ConnectClient(settings, tinkoffInvestBot.Update.FromChat().ID, tinkoffInvestBot.Update.Message.Text)
				if client != nil {
					mainMenu(tinkoffInvestBot)
				}
			}
		}
	}
}

func mainMenu(tinkoffInvestBot *model.Bot) {
	msg := tgbotapi.NewMessage(tinkoffInvestBot.Update.FromChat().ID, "Выберите что хотите сделать: \n"+
		"0 - Закончить\n"+
		"1 - Посмотреть свои ценные бумаги\n"+
		"2 - Посмотреть данные по акции по ее тикеру\n"+
		"3 - Добавить акцию не из вашего портфеля в список отслеживаемых\n"+
		"4 - Посмотреть расписание дивидендов по вашим ценным бумагам\n"+
		"5 - Посмотреть расписание отчётов по вашим ценным бумагам")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("0", "command=0"),
			tgbotapi.NewInlineKeyboardButtonData("1", "command=1"),
			tgbotapi.NewInlineKeyboardButtonData("2", "command=2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("3", "command=3"),
			tgbotapi.NewInlineKeyboardButtonData("4", "command=4"),
			tgbotapi.NewInlineKeyboardButtonData("5", "command=5"),
		),
	)
	if _, err := tinkoffInvestBot.Api.Send(msg); err != nil {
		log.Fatalln(err.Error())
	}
}

func printResponseToBot(tinkoffInvestBot *model.Bot, responses []string, logger *zap.SugaredLogger, numericKeyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(tinkoffInvestBot.Update.FromChat().ID, "Данные отсутствуют")
	if responses != nil {
		for _, response := range responses {
			msg.Text = response
			if numericKeyboard != nil {
				msg.ReplyMarkup = numericKeyboard
			}
			if _, err := tinkoffInvestBot.Api.Send(msg); err != nil {
				logger.Fatalln(err.Error())
			}
		}
	} else {
		if _, err := tinkoffInvestBot.Api.Send(msg); err != nil {
			logger.Fatalln(err.Error())
		}
	}
}
