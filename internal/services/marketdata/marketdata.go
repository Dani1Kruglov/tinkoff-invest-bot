package marketdata

import (
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	"tinkoff-investment-bot/internal/model"
	printbot "tinkoff-investment-bot/internal/print-bot"
)

func GetLastPriceByFigi(tracker *model.Tracker, instrument *investapi.Share, logger *zap.SugaredLogger) {
	marketDataResp, err := tracker.MarketDataService.GetLastPrices([]string{instrument.GetFigi()})
	if err != nil {
		logger.Errorf(err.Error())
	}
	printbot.LastPrice(marketDataResp.GetLastPrices()[0].GetPrice().GetUnits(), marketDataResp.GetLastPrices()[0].GetPrice().GetNano()/10000000)
}
