package users

import (
	"fmt"
	"tinkoff-investment-bot/internal/model"
)

func GetAccountID(tracker *model.Tracker) (string, error) {
	userInfo, err := tracker.UsersService.GetAccounts()
	if err != nil {
		return "", err
	}

	accounts := userInfo.GetAccounts()
	if len(userInfo.GetAccounts()) > 1 {
		var num int
		fmt.Println("Выберите счёт, который хотите просмотреть: ")
		for i, account := range accounts {
			fmt.Println("-  -  -  -  -  -  -  -  -  -  -")
			fmt.Printf("%d: \n id счёта: %v\n", i+1, account.GetId())
			fmt.Printf("id счёта: %v\n", account.GetName())
		}
		_, err = fmt.Scan(&num)
		if err != nil {
			return "", err
		}
		return accounts[num].GetId(), nil

	}
	return accounts[0].GetId(), nil
}
