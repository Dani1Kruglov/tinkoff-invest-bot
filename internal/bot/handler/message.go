package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"log"
	"tinkoff-investment-bot/internal/bot/model"
	"tinkoff-investment-bot/internal/connect/client"
	ms "tinkoff-investment-bot/internal/model/settings"
	"tinkoff-investment-bot/internal/tracking"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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

func MessageHandler(tinkoffInvestBot *model.Bot, settings *ms.Settings, clientInvest *investgo.Client, cancel *context.CancelFunc) {
	if tinkoffInvestBot.Update.Message.Text[:5] == "token" {
		clientInvest, *cancel = client.ConnectClient(settings, tinkoffInvestBot.Update.FromChat().ID, tinkoffInvestBot.Update.Message.Text[6:])
	} else if tinkoffInvestBot.Update.Message.Text[:6] == "ticker" && clientInvest != nil {
		responses := tracking.GetShare(settings, clientInvest, tinkoffInvestBot.Update.Message.Text[7:])
		printRespToBot(tinkoffInvestBot, responses)
	}
	if clientInvest != nil {
		MainMenu(tinkoffInvestBot)
	}

}

func MainMenu(tinkoffInvestBot *model.Bot) {
	msg := tgbotapi.NewMessage(tinkoffInvestBot.Update.FromChat().ID, "Выберите что хотите сделать: \n"+
		"0 - Закончить\n"+
		"1 - Посмотреть свои ценные бумаги\n"+
		"2 - Посмотреть данные по акции по ее тикеру\n"+
		"3 - Добавить акцию не из вашего портфеля в список отслеживаемых\n"+
		"4 - Посмотреть расписание дивидендов по вашим ценным бумагам\n"+
		"5 - Посмотреть расписание отчётов по вашим ценным бумагам")
	msg.ReplyMarkup = numericKeyboard
	if _, err := tinkoffInvestBot.Api.Send(msg); err != nil {
		log.Fatalln(err.Error())
	}
}

func printRespToBot(tinkoffInvestBot *model.Bot, responses []string) {
	msg := tgbotapi.NewMessage(tinkoffInvestBot.Update.FromChat().ID, "Данные отсутствуют")
	if responses != nil {
		for _, response := range responses {
			msg.Text = response
			if _, err := tinkoffInvestBot.Api.Send(msg); err != nil {
				log.Fatalln(err.Error())
			}
		}
	} else {
		if _, err := tinkoffInvestBot.Api.Send(msg); err != nil {
			log.Fatalln(err.Error())
		}
	}
}
