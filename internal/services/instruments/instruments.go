package instruments

import (
	"fmt"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"go.uber.org/zap"
)

func WorkWithInstruments(client *investgo.Client, logger *zap.SugaredLogger) {

	instrumentsService := client.NewInstrumentsServiceClient()
	marketDataService := client.NewMarketDataServiceClient()

	instrResp, err := instrumentsService.ShareByTicker("MOEX", "TQBR")

	if err != nil {
		logger.Errorf(err.Error())
	} else {
		instrument := instrResp.GetInstrument()
		fmt.Println(instrument)
		fmt.Printf("-------\n")
		forecast, _ := instrumentsService.GetForecastBy(instrument.GetUid())
		for _, target := range forecast.Targets {
			fmt.Println(target.TargetPrice)
			fmt.Println(target.CurrentPrice)
			fmt.Println(target.PriceChange)
		}
		fmt.Println("/////////")
		fmt.Println(forecast.Consensus.Consensus)

		fmt.Printf("++++++++\n")
		marketDataResp, _ := marketDataService.GetLastPrices([]string{instrument.GetFigi()})
		fmt.Println(marketDataResp)
	}
}
