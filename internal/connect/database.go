package connect

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"tinkoff-investment-bot/internal/config"
	"tinkoff-investment-bot/internal/model"
)

func connectDB() (*gorm.DB, error) {
	databaseDSN, err := config.LoadConfigDBFile("config.yaml")
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(postgres.Open(databaseDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)

	err = db.AutoMigrate(&model.User{}, &model.Account{}, &model.Share{}, &model.UserFavoriteShare{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
