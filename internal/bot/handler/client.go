package handler

import (
	"context"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"tinkoff-investment-bot/internal/connect/client"
	ms "tinkoff-investment-bot/internal/model/settings"
	"tinkoff-investment-bot/internal/storage"
)

func ClientHandler(telegramChatID int64, settings *ms.Settings) (*investgo.Client, context.CancelFunc) {
	userStorage := storage.NewUserStorage(settings.DB)
	user := userStorage.GetUserByTelegramChatID(telegramChatID)

	if user.ID != 0 {
		clientInvest, cancel := client.ConnectClient(settings, telegramChatID, user.Token)
		return clientInvest, cancel
	}

	return nil, nil
}
