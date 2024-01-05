package instruments

import (
	"fmt"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	"sync"
	"time"
	"tinkoff-investment-bot/internal/model"
	printbot "tinkoff-investment-bot/internal/print-bot"
	m "tinkoff-investment-bot/internal/services/marketdata"
	o "tinkoff-investment-bot/internal/services/operations"
)

func GetShareByTicker(tracker *model.Tracker, logger *zap.SugaredLogger) {
	ticker, err := printbot.GetTickerFromUser()
	if err != nil {
		logger.Errorf(err.Error())
	}

	instrResp, err := tracker.InstrumentsService.ShareByTicker(ticker, "TQBR") //tqbr только для российских компаний

	if err != nil {
		logger.Errorf(err.Error())
	} else {
		instrument := instrResp.GetInstrument()
		printbot.InfoAboutShareByItsTicker(instrument)

		m.GetLastPriceByFigi(tracker, instrument, logger)

		printbot.HeadlineForecastsOfInvestmentHouses()

		forecast, _ := tracker.InstrumentsService.GetForecastBy(instrument.GetUid())
		for i, target := range forecast.GetTargets() {
			printbot.InvestHouseForecast(i, target)
		}
		fmt.Println("======================================")

		printbot.ConsensusForecast(forecast.GetConsensus().GetConsensus().GetUnits(), forecast.GetConsensus().GetConsensus().GetNano()/10000000)
	}
	fmt.Println("[][][][[][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][]")

}

func GetScheduleOnClientSecurities(tracker *model.Tracker, logger *zap.SugaredLogger, isReports bool) {
	portfolioResp, err := o.GetPortfolioByAccountID(tracker)

	if err != nil {
		logger.Errorf(err.Error())
	} else {
		var wg sync.WaitGroup
		for _, position := range portfolioResp.GetPositions() {
			fmt.Println("-  -  -  -  -  -  -  -  -  -  -")
			instrument, err := tracker.InstrumentsService.InstrumentByUid(position.GetInstrumentUid())
			if err != nil {
				logger.Errorf(err.Error())
			}
			wg.Add(1)
			position := position
			go func() {
				defer wg.Done()
				err := getPaperWithShareTypeFromInstruments(instrument.GetInstrument(), tracker, position, isReports)
				if err != nil {
					logger.Errorf(err.Error())
				}
			}()

		}
		wg.Wait()
	}
}

func getPaperWithShareTypeFromInstruments(instrument *investapi.Instrument, tracker *model.Tracker, position *investapi.PortfolioPosition, isReports bool) error {
	if instrument.GetInstrumentType() != "share" { //==
		err := getSchedule(tracker, instrument, position, isReports)
		if err != nil {
			return err
		}
	}
	return nil
}

func getSchedule(tracker *model.Tracker, instrument *investapi.Instrument, position *investapi.PortfolioPosition, isReports bool) error {
	startDate := time.Now()
	endDate := startDate.AddDate(0, 6, 0)

	//ticker, _ := tracker.InstrumentsService.ShareByTicker("TATNP", "TQBR")

	if isReports {
		reports, err := tracker.InstrumentsService.GetAssetReports(position.GetInstrumentUid(), startDate, endDate)
		if err != nil {
			return err
		}

		printbot.PrintScheduleOfReports(instrument.GetName(), reports)
	} else {
		reports, err := tracker.InstrumentsService.GetDividents(position.GetInstrumentUid(), startDate, endDate)
		if err != nil {
			return err
		}
		printbot.PrintScheduleOfDividend(instrument.GetName(), reports.GetDividends()[0].GetPaymentDate().AsTime().Local(), reports.GetDividends()[0].GetDividendNet())
	}
	return nil
}
