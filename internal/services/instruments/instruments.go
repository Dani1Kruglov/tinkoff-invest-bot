package instruments

import (
	"fmt"
	"go.uber.org/zap"

	"tinkoff-investment-bot/internal/model"
	m "tinkoff-investment-bot/internal/services/marketdata"
)

func GetShareByTicker(logger *zap.SugaredLogger, tracker *model.Tracker) {
	fmt.Println("Введите тикер акции (MOEX, SBER и так далее)")
	var ticker string
	_, err := fmt.Scan(&ticker)
	if err != nil {
		logger.Errorf(err.Error())
	}

	instrResp, err := tracker.InstrumentsService.ShareByTicker(ticker, "TQBR")

	if err != nil {
		logger.Errorf(err.Error())
	} else {
		instrument := instrResp.GetInstrument()
		fmt.Println(instrument)

		marketDataResp, err := m.GetLastPriceByFigi(instrument, tracker)
		if err != nil {
			logger.Errorf(err.Error())
		}
		fmt.Printf("Нынешняя цена: %v.%v ₽\n", marketDataResp.GetLastPrices()[0].GetPrice().GetUnits(), marketDataResp.GetLastPrices()[0].GetPrice().GetNano()/10000000)

		fmt.Println("Прогнозы от инвест домов:")
		forecast, _ := tracker.InstrumentsService.GetForecastBy(instrument.GetUid())
		for i, target := range forecast.GetTargets() {
			fmt.Println("-  -  -  -  -  -  -  -  -  -  -")
			fmt.Printf("Номер компании: %d\n", i+1)
			fmt.Printf("Название компании, давшей прогноз: %v\n", target.GetCompany())
			fmt.Printf("Прогноз: %v\n", target.GetRecommendation())
			fmt.Printf("Прогнозируемая цена: %v.%v ₽\n", target.GetTargetPrice().GetUnits(), target.GetTargetPrice().GetNano()/10000000)
			fmt.Printf("Изменение цены: %v.%v ₽\n", target.GetPriceChange().GetUnits(), target.GetPriceChange().GetNano()/10000000)
		}
		fmt.Println("======================================")
		fmt.Printf("Согласованный прогноз: %v.%v ₽\n", forecast.GetConsensus().GetConsensus().GetUnits(), forecast.GetConsensus().GetConsensus().GetNano()/10000000)
	}
	fmt.Println("[][][][[][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][]")

}
