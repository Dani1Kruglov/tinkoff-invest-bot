package tracking

import (
	"fmt"
	"log"
	printbot "tinkoff-investment-bot/internal/print-bot"
	i "tinkoff-investment-bot/internal/services/instruments"

	"tinkoff-investment-bot/internal/connect"
	"tinkoff-investment-bot/internal/model"

	o "tinkoff-investment-bot/internal/services/operations"
	u "tinkoff-investment-bot/internal/services/users"
)

func TrackByTinkoffToken() {
	token, _ := printbot.GetTokenFromUser()

	command := printbot.MainMenu()

	client, logger, cancel, _ := connect.ClientByConfig(token)

	var tracker model.Tracker
	tracker.AddServices(client)

	end := false
	for !end {
		switch command {
		case "0":
			end = true
			break
		case "1":
			o.GetUserSecuritiesOnAccount(&tracker, logger)
			break
		case "2":
			i.GetShareByTicker(&tracker, logger)
			break
		case "3":
			accountID, err := u.GetAccountID(&tracker)
			if err != nil {
				logger.Errorf(err.Error())
			}
			fmt.Println(accountID)
			fmt.Println("[][][][[][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][]")
			break
		case "4":
			i.GetScheduleOnClientSecurities(&tracker, logger, false)
			break
		case "5":
			i.GetScheduleOnClientSecurities(&tracker, logger, true)
			break
		default:
			break

		}
		command = printbot.MainMenu()
	}

	defer func() {
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

	}()
}
