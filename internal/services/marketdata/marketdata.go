package marketdata

import (
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"tinkoff-investment-bot/internal/model"
)

func GetLastPriceByFigi(tracker *model.Tracker, instrument *investapi.Share) (float32, error) {
	marketDataResp, err := tracker.MarketDataService.GetLastPrices([]string{instrument.GetFigi()})
	if err != nil {
		return 0, err
	}
	return float32(marketDataResp.GetLastPrices()[0].GetPrice().GetUnits()) + float32(marketDataResp.GetLastPrices()[0].GetPrice().GetNano())/float32(1000000000), nil
}
