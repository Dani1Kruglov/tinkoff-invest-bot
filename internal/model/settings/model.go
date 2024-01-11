package settings

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	db "tinkoff-investment-bot/internal/connect/database"
	l "tinkoff-investment-bot/internal/connect/logger"
)

type Settings struct {
	Logger *zap.SugaredLogger
	DB     *gorm.DB
}

func NewSettings() Settings {
	logger := l.GetLogger()
	return Settings{
		Logger: logger,
		DB:     db.ConnectDB(logger),
	}
}
