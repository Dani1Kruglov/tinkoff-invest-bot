package tracking

import (
	printbot "tinkoff-investment-bot/internal/print-bot"
	"tinkoff-investment-bot/internal/services/instruments/schedules"
	"tinkoff-investment-bot/internal/services/instruments/shares"

	"tinkoff-investment-bot/internal/connect"

	o "tinkoff-investment-bot/internal/services/operations"
)

func TrackByTinkoffToken() {
	telegramID := "telegramID" //get from telegram bot

	client, db, cancel, logger, tracker := connect.Config(telegramID)
	defer func() {
		connect.Close(client, db, cancel, logger)
	}()

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
			schedules.GetScheduleOnClientSecurities(tracker, logger, db, telegramID, false)
			break
		case "5":
			schedules.GetScheduleOnClientSecurities(tracker, logger, db, telegramID, true)
			break
		default:
			break
		}
	}
}
