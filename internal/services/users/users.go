package users

import (
	"tinkoff-investment-bot/internal/model"
	printbot "tinkoff-investment-bot/internal/print-bot"
)

func GetAccount(tracker *model.Tracker) (model.Account, error) {
	userInfo, err := tracker.UsersService.GetAccounts()
	if err != nil {
		return model.Account{}, err
	}

	accounts := userInfo.GetAccounts()

	if len(userInfo.GetAccounts()) > 1 {
		accountSelect, err := printbot.UserAccountSelect(accounts)
		if err != nil {
			return model.Account{}, err
		}
		account := model.Account{
			AccountID: accounts[accountSelect].GetId(),
			Name:      accounts[accountSelect].GetName(),
		}
		return account, nil

	}
	account := model.Account{
		AccountID: accounts[0].GetId(),
		Name:      accounts[0].GetName(),
	}
	return account, nil
}
