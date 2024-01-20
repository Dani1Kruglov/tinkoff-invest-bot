package listener

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
	"tinkoff-investment-bot/internal/bot/checker"
	"tinkoff-investment-bot/internal/bot/model"
	ms "tinkoff-investment-bot/internal/model/settings"
)

var (
	cacheCommand sync.Map
)

func ListenUpdates(tinkoffInvestBot *model.Bot, settings *ms.Settings) {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := tinkoffInvestBot.Api.GetUpdatesChan(u)

	for {
		select {
		case tinkoffInvestBot.Update = <-updates:
			checker.CheckUpdate(tinkoffInvestBot, settings, &cacheCommand)
		}
	}
}
