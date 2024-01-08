package config

import (
	"context"
	"fmt"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
	printbot "tinkoff-investment-bot/internal/print-bot"
)

func ClientTinkoffInvestByConfigYaml(logger *zap.SugaredLogger, token *string) (*investgo.Client, context.CancelFunc, error) {
	config, err := investgo.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("config loading error %v", err.Error())
	}

	config.Token = *token

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	client, err := investgo.NewClient(ctx, config, logger)
	for err != nil {
		if err.Error() == "rpc error: code = Unauthenticated desc = 40003" {
			fmt.Println("Неверный токен, попробуйте еще раз: ")
			*token, err = printbot.GetTokenFromUser()
			config.Token = *token
			client, err = investgo.NewClient(ctx, config, logger)
		} else {
			logger.Fatalf("connect creating error %v", err.Error())
		}
	}

	return client, cancel, nil
}
