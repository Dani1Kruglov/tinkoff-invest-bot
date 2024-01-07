package tracking

import (
	"fmt"
	printbot "tinkoff-investment-bot/internal/print-bot"
	i "tinkoff-investment-bot/internal/services/instruments"

	"tinkoff-investment-bot/internal/connect"

	o "tinkoff-investment-bot/internal/services/operations"
	u "tinkoff-investment-bot/internal/services/users"
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
			i.GetShareByTicker(tracker, logger)
			break
		case "3":
			account, err := u.GetAccount(tracker)
			if err != nil {
				logger.Errorf(err.Error())
			}
			fmt.Println(account.UserID)
			fmt.Println("[][][][[][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][]")
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
