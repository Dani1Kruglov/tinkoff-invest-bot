package users

import (
	printbot "tinkoff-investment-bot/internal/bot/print"
	"tinkoff-investment-bot/internal/model/database"
	"tinkoff-investment-bot/internal/model/tracker"
)

func GetAccount(tracker *tracker.Tracker) (database.Account, error) {
	userInfo, err := tracker.UsersService.GetAccounts()
	if err != nil {
		return database.Account{}, err
	}

	accounts := userInfo.GetAccounts()

	if len(userInfo.GetAccounts()) > 1 {
		accountSelect, err := printbot.UserAccountSelect(accounts)
		if err != nil {
			return database.Account{}, err
		}
		account := database.Account{
			AccountID: accounts[accountSelect].GetId(),
			Name:      accounts[accountSelect].GetName(),
		}
		return account, nil

	}
	account := database.Account{
		AccountID: accounts[0].GetId(),
		Name:      accounts[0].GetName(),
	}
	return account, nil
}
