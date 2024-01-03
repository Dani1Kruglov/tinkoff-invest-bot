package main

import (
	"fmt"
	"log"
	"tinkoff-investment-bot/internal/connect"
	i "tinkoff-investment-bot/internal/services/instruments"
	o "tinkoff-investment-bot/internal/services/operations"
)

func main() {
	client, logger, cancel, _ := connect.ClientByConfig()

	i.WorkWithInstruments(client, logger)
	o.WorkWithOperations(client, logger)

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

/* get account data
usersService := client.NewUsersServiceClient()

accsResp, err := usersService.GetAccounts()
if err != nil {
logger.Errorf(err.Error())
} else {
accs := accsResp.GetAccounts()
fmt.Println(accs)
}*/
