package schedules

import (
	"fmt"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"time"
	"tinkoff-investment-bot/internal/model"
	printbot "tinkoff-investment-bot/internal/print-bot"
	o "tinkoff-investment-bot/internal/services/operations"
)

func GetScheduleOnClientSecurities(tracker *model.Tracker, logger *zap.SugaredLogger, db *gorm.DB, telegramID string, isReports bool) {
	portfolioResp, err := o.GetPortfolioByAccountID(tracker, db, telegramID)

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
	if instrument.GetInstrumentType() == "share" {
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
