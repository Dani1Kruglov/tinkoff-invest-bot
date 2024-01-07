package connect

import (
	"context"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"log"
	"os/signal"
	"syscall"
	"time"
	"tinkoff-investment-bot/internal/database"
	"tinkoff-investment-bot/internal/model"
	printbot "tinkoff-investment-bot/internal/print-bot"
)

func Config(telegramID string) (*investgo.Client, *gorm.DB, context.CancelFunc, *zap.SugaredLogger, *model.Tracker) {
	logger := getLogger()

	db, err := connectDB()
	if err != nil {
		logger.Errorf(err.Error())
	}

	user := database.GetUserByTelegramID(db, telegramID)

	var token string

	if user.ID != 0 {
		token = user.Token
	} else {
		token, err = printbot.GetTokenFromUser()
		if err != nil {
			logger.Errorf(err.Error())
		}
		err := database.AddUser(db, &model.User{TelegramID: telegramID, Token: token})
		if err != nil {
			logger.Errorf(err.Error())
		}
	}

	client, cancel := clientByConfig(logger, token)

	var tracker model.Tracker
	tracker.AddServices(client)

	return client, db, cancel, logger, &tracker
}

func getLogger() *zap.SugaredLogger {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	zapConfig.EncoderConfig.TimeKey = "time"
	l, err := zapConfig.Build()
	if err != nil {
		log.Fatalf("logger creating error %v", err)
	}
	return l.Sugar()
}

func clientByConfig(logger *zap.SugaredLogger, token string) (*investgo.Client, context.CancelFunc) {
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

	return client, cancel
}
