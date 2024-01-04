package main

import (
	"fmt"
	"tinkoff-investment-bot/internal/tracking"
)

func main() {

	fmt.Println("Хотите использовать бота для: ")
	fmt.Println("1 - только для отслеживания информации по ценным бумагам ")
	fmt.Println("2 - для выполнения операций ")
	var command string
	_, err := fmt.Scan(&command)
	if err != nil {
		_ = fmt.Errorf(err.Error())
	}
	switch command {
	case "1":
		tracking.TrackByTinkoffToken()
		break
	case "2":
		break
	default:
		fmt.Println("Такой команды нет")
		break
	}
}
