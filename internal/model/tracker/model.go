package tracker

import "github.com/russianinvestments/invest-api-go-sdk/investgo"

type Tracker struct {
	InstrumentsService *investgo.InstrumentsServiceClient
	MarketDataService  *investgo.MarketDataServiceClient
	OperationsService  *investgo.OperationsServiceClient
	UsersService       *investgo.UsersServiceClient
}

func NewTracker(client *investgo.Client) *Tracker {
	return &Tracker{
		InstrumentsService: client.NewInstrumentsServiceClient(),
		MarketDataService:  client.NewMarketDataServiceClient(),
		OperationsService:  client.NewOperationsServiceClient(),
		UsersService:       client.NewUsersServiceClient(),
	}
}
