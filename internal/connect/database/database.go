package database

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"tinkoff-investment-bot/internal/connect/config"
	"tinkoff-investment-bot/internal/model/database"
)

func ConnectDB(logger *zap.SugaredLogger) *gorm.DB {
	databaseDSN, err := config.LoadConfigDBFileByConfigYaml("config.yaml")
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	db, err := gorm.Open(postgres.Open(databaseDSN), &gorm.Config{})
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)

	err = db.AutoMigrate(&database.User{}, &database.Account{}, &database.Share{}, &database.UserFavoriteShare{})
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	return db
}
