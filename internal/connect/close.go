package connect

import (
	"log"
	s "tinkoff-investment-bot/internal/model/settings"
)

func Close(settings *s.Settings) {
	settings.Logger.Infoln("-------end-------")
	err := settings.Logger.Sync()
	if err != nil {
		log.Printf(err.Error())
	}

	settings.Logger.Infof("closing connect connection")

	settings.Logger.Infof("closing database connection")
	sqlDB, err := settings.DB.DB()
	err = sqlDB.Close()
	if err != nil {
		settings.Logger.Errorf("close db connect error %v", err.Error())
	}
}
