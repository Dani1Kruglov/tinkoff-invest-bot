package tinkoff_client

import (
	"context"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"tinkoff-investment-bot/internal/connect/config"
	"tinkoff-investment-bot/internal/model/database"
	"tinkoff-investment-bot/internal/model/settings"
	"tinkoff-investment-bot/internal/storage"
)

func ConnectClient(settings *settings.Settings, telegramChatID int64, token string) (*investgo.Client, context.CancelFunc) {

	client, cancel, err := config.ClientTinkoffInvestByConfigYaml(settings.Logger, token)
	if err != nil {
		settings.Logger.Errorf(err.Error())
	}

	userStorage := storage.NewUserStorage(settings.DB)
	err = userStorage.AddUser(&database.User{TelegramID: telegramChatID,
		Token: token})
	if err != nil {
		settings.Logger.Errorf(err.Error())
	}

	return client, cancel
}
