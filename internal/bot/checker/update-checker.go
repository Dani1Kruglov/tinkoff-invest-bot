package checker

import (
	"sync"
	"tinkoff-investment-bot/internal/bot/handler"
	"tinkoff-investment-bot/internal/bot/model"
	ms "tinkoff-investment-bot/internal/model/settings"
)

func CheckUpdate(tinkoffInvestBot *model.Bot, settings *ms.Settings, cacheCommand *sync.Map) {

	client, cancel := handler.ClientHandler(tinkoffInvestBot.Update.FromChat().ID, settings)

	defer func() {
		if client != nil {
			err := client.Stop()
			if err != nil {
				settings.Logger.Errorf("connect shutdown error %v", err.Error())
			}
		}
		if cancel != nil {
			cancel()
		}
	}()

	if tinkoffInvestBot.Update.Message != nil {
		if tinkoffInvestBot.Update.Message.Command() != "" {
			handler.CommandHandler(tinkoffInvestBot)
		} else if tinkoffInvestBot.Update.Message.Text != "" {
			handler.MessageHandler(tinkoffInvestBot, settings, client, &cancel, cacheCommand)
		}
	} else if tinkoffInvestBot.Update.CallbackQuery != nil {
		handler.InlineKeyBoardHandler(tinkoffInvestBot, settings, client, cacheCommand)
	}
}
