package tracking

import (
	"fmt"
	"log"
	i "tinkoff-investment-bot/internal/services/instruments"

	"tinkoff-investment-bot/internal/connect"
	"tinkoff-investment-bot/internal/model"

	o "tinkoff-investment-bot/internal/services/operations"
	u "tinkoff-investment-bot/internal/services/users"
)

func TrackByTinkoffToken() {
	fmt.Println("Введите 'токен для чтения' своего брокерского счета:")
	var token, command string
	_, err := fmt.Scan(&token)
	if err != nil {
		_ = fmt.Errorf(err.Error())
	}
	client, logger, cancel, _ := connect.ClientByConfig(token)

	var tracker model.Tracker
	tracker.AddServices(client)

	end := false
	for !end {
		fmt.Println("Введите что хотите сделать: ")
		fmt.Println("0 - Закончить")
		fmt.Println("1 - Посмотреть свои ценные бумаги")
		fmt.Println("2 - Посмотреть данные по акции по ее названию")
		fmt.Println("3 - Добавить акцию не из вашего портфеля в список отслеживаемых ")
		fmt.Println("4 - Посмотреть расписание дивидендов по вашим ценным бумагам")
		fmt.Println("5 - Посмотреть расписание отчётов по вашим ценным бумагам")
		_, err := fmt.Scan(&command)
		if err != nil {
			_ = fmt.Errorf(err.Error())
		}
		switch command {
		case "0":
			end = true
			break
		case "1":
			o.GetUserSharesOnAccount(logger, &tracker)
			break
		case "2":
			i.GetShareByTicker(logger, &tracker)
			break
		case "3":
			accountID, err := u.GetAccountID(&tracker)
			if err != nil {
				logger.Errorf(err.Error())
			}
			fmt.Println(accountID)
			fmt.Println("[][][][[][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][]")
			break
		default:
			break

		}
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
