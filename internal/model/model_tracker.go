package model

import "github.com/russianinvestments/invest-api-go-sdk/investgo"

type Tracker struct {
	InstrumentsService *investgo.InstrumentsServiceClient
	MarketDataService  *investgo.MarketDataServiceClient
	OperationsService  *investgo.OperationsServiceClient
	UsersService       *investgo.UsersServiceClient
}

func (t *Tracker) AddServices(client *investgo.Client) {
	t.InstrumentsService = client.NewInstrumentsServiceClient()
	t.MarketDataService = client.NewMarketDataServiceClient()
	t.OperationsService = client.NewOperationsServiceClient()
	t.UsersService = client.NewUsersServiceClient()
}
