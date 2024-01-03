package connect

import (
	"context"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func ClientByConfig() (*investgo.Client, *zap.SugaredLogger, context.CancelFunc, investgo.Config) {
	config, err := investgo.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("config loading error %v", err.Error())
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	zapConfig.EncoderConfig.TimeKey = "time"
	l, err := zapConfig.Build()
	logger := l.Sugar()

	if err != nil {
		log.Fatalf("logger creating error %v", err)
	}

	client, err := investgo.NewClient(ctx, config, logger)
	if err != nil {
		logger.Fatalf("connect creating error %v", err.Error())
	}

	return client, logger, cancel, config
}
