package marketdata

import (
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"tinkoff-investment-bot/internal/model"
)

func GetLastPriceByFigi(instrument *investapi.Share, tracker *model.Tracker) (*investgo.GetLastPricesResponse, error) {
	return tracker.MarketDataService.GetLastPrices([]string{instrument.GetFigi()})
}
