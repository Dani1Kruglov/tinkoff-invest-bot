package invest_schedules

import (
	"fmt"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"time"
	"tinkoff-investment-bot/internal/model/tracker"
	o "tinkoff-investment-bot/internal/services/operations"
)

func GetScheduleOnClientSecurities(tracker *tracker.Tracker, logger *zap.SugaredLogger, db *gorm.DB, telegramChatID int64, isReports bool) []string {
	portfolioResp, err := o.GetPortfolioByAccountID(tracker, db, telegramChatID)
	var schedules []string

	if err != nil {
		logger.Errorf(err.Error())
	} else {
		var wg sync.WaitGroup
		for _, position := range portfolioResp.GetPositions() {
			instrument, err := tracker.InstrumentsService.InstrumentByUid(position.GetInstrumentUid())
			if err != nil {
				logger.Errorf(err.Error())
			}
			wg.Add(1)
			position := position
			go func() {
				defer wg.Done()
				shareSchedule, err := getPaperWithShareTypeFromInstruments(instrument.GetInstrument(), tracker, position, isReports)
				if err != nil {
					logger.Errorf(err.Error())
				}
				if shareSchedule != "" {
					schedules = append(schedules, shareSchedule)
				}
			}()

		}
		wg.Wait()
	}
	return schedules
}

func getPaperWithShareTypeFromInstruments(instrument *investapi.Instrument, tracker *tracker.Tracker, position *investapi.PortfolioPosition, isReports bool) (string, error) {
	if instrument.GetInstrumentType() == "share" {
		share, err := getSchedule(tracker, instrument, position, isReports)
		if err != nil {
			return "", err
		}
		return share, err
	}
	return "", nil
}

func getSchedule(tracker *tracker.Tracker, instrument *investapi.Instrument, position *investapi.PortfolioPosition, isReports bool) (string, error) {
	startDate := time.Now()
	endDate := startDate.AddDate(0, 6, 0)

	if isReports {
		reports, err := tracker.InstrumentsService.GetAssetReports(position.GetInstrumentUid(), startDate, endDate)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Отчет по бумаге %v в ближайшие 6 месяцев: %v\n", instrument.GetName(), reports.GetAssetReportsResponse), nil

	}
	reports, err := tracker.InstrumentsService.GetDividents(position.GetInstrumentUid(), startDate, endDate)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Дивиденды по бумаге '%v' в ближайшие 6 месяцев:\nДата фактических выплат по МСК: %v\nВеличина дивиденда на 1 ценную бумагу (включая валюту): %d.%d\n",
		instrument.GetName(), reports.GetDividends()[0].GetPaymentDate().AsTime().Local(),
		reports.GetDividends()[0].GetDividendNet().GetUnits(),
		reports.GetDividends()[0].GetDividendNet().GetNano()), nil
}
