package users

import (
	"tinkoff-investment-bot/internal/model"
	printbot "tinkoff-investment-bot/internal/print-bot"
)

func GetAccountID(tracker *model.Tracker) (string, error) {
	userInfo, err := tracker.UsersService.GetAccounts()
	if err != nil {
		return "", err
	}

	accounts := userInfo.GetAccounts()

	if len(userInfo.GetAccounts()) > 1 {
		accountSelect, err := printbot.UserAccountSelect(accounts)
		if err != nil {
			return "", err
		}
		return accounts[accountSelect].GetId(), nil

	}
	return accounts[0].GetId(), nil
}
