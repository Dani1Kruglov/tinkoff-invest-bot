package connect

import (
	"context"
	"fmt"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
)

func Close(client *investgo.Client, db *gorm.DB, cancel context.CancelFunc, logger *zap.SugaredLogger) {
	fmt.Println("-------end-------")
	err := logger.Sync()
	if err != nil {
		log.Printf(err.Error())
	}

	cancel()

	logger.Infof("closing connect connection")
	err = client.Stop()
	if err != nil {
		logger.Errorf("connect shutdown error %v", err.Error())
	}

	logger.Infof("closing database connection")
	sqlDB, err := db.DB()
	err = sqlDB.Close()
	if err != nil {
		logger.Errorf("close db connect error %v", err.Error())
	}
}
