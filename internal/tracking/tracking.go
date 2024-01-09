package tracking

import (
	"tinkoff-investment-bot/internal/connect"
	"tinkoff-investment-bot/internal/cron"
	printbot "tinkoff-investment-bot/internal/print-bot"
	is "tinkoff-investment-bot/internal/services/instruments/invest-schedules"
	"tinkoff-investment-bot/internal/services/instruments/shares"

	o "tinkoff-investment-bot/internal/services/operations"
)

func TrackByTinkoffToken() {
	telegramID := "telegramID" //get from telegram bot

	client, db, cancel, logger, tracker := connect.Config(telegramID)
	defer func() {
		connect.Close(client, db, cancel, logger)
	}()

	logger.Infoln("start cron schedule")
	go cron.NewCron(db, tracker)

	end := false
	for !end {
		command := printbot.MainMenu()
		switch command {
		case "0":
			end = true
			break
		case "1":
			o.GetUserSecuritiesOnAccount(tracker, logger, db, telegramID)
			break
		case "2":
			shares.ViewInfoOnShareByItsTicker(tracker, logger)
			break
		case "3":
			shares.AddShareToListOfTracked(tracker, logger, db, telegramID)
			break
		case "4":
			is.GetScheduleOnClientSecurities(tracker, logger, db, telegramID, false)
			break
		case "5":
			is.GetScheduleOnClientSecurities(tracker, logger, db, telegramID, true)
			break
		default:
			break
		}
	}
}
