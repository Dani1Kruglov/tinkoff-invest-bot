package tracking

import (
	printbot "tinkoff-investment-bot/internal/print-bot"
	i "tinkoff-investment-bot/internal/services/instruments"

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
			i.ViewInfoOnShareByItsTicker(tracker, logger)
			break
		case "3":
			i.AddShareToListOfTracked(tracker, logger, db, telegramID)
			break
		case "4":
			i.GetScheduleOnClientSecurities(tracker, logger, db, telegramID, false)
			break
		case "5":
			i.GetScheduleOnClientSecurities(tracker, logger, db, telegramID, true)
			break
		default:
			break
		}
	}
}
