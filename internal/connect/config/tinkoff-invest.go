package config

import (
	"context"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
)

func ClientTinkoffInvestByConfigYaml(logger *zap.SugaredLogger, token string) (*investgo.Client, context.CancelFunc, error) {
	config, err := investgo.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("config loading error %v", err.Error())
	}

	config.Token = token

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	client, err := investgo.NewClient(ctx, config, logger)

	if err != nil {
		logger.Fatalf("connect creating error %v", err.Error())
	}

	return client, cancel, nil
}
