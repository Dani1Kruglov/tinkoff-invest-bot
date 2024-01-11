package shares

import (
	"fmt"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	printbot "tinkoff-investment-bot/internal/bot/print"
	"tinkoff-investment-bot/internal/model/database"
	"tinkoff-investment-bot/internal/model/tracker"
	m "tinkoff-investment-bot/internal/services/marketdata"
	"tinkoff-investment-bot/internal/storage"
)

func ViewInfoOnShareByItsTicker(tracker *tracker.Tracker, logger *zap.SugaredLogger, ticker string) []string {
	instrument, shareResp, _, err := getInfoAboutShareByTicker(tracker, ticker)
	if err != nil {
		logger.Errorf(err.Error())
	}
	shareResp = append(shareResp, getForecastsAboutShare(tracker, instrument)...)
	return shareResp
}

func AddShareToListOfTracked(tracker *tracker.Tracker, logger *zap.SugaredLogger, db *gorm.DB, telegramChatID int64) {
	instrument, _, price32, err := getInfoAboutShareByTicker(tracker, "")
	if err != nil {
		logger.Errorf(err.Error())
	}

	command, err := printbot.AddToListOfTracked()
	if err != nil {
		logger.Errorf(err.Error())
	}

	if command == "1" {
		var price float64
		command, err = printbot.SpecifyPrice()
		if err != nil {
			logger.Errorf(err.Error())
		}

		if command != "1" {
			price, err = strconv.ParseFloat(command, 32)
			if err != nil {
				logger.Errorf(err.Error())
			}
			price32 = float32(price)
		}

		userStorage := storage.NewUserStorage(db)
		shareStorage := storage.NewShareStorage(db)

		err = shareStorage.AddShare(&database.Share{
			UID:       instrument.GetUid(),
			Ticker:    instrument.GetTicker(),
			Name:      instrument.GetName(),
			FIGI:      instrument.GetFigi(),
			ClassCode: instrument.GetClassCode(),
		}, userStorage.GetUserByTelegramChatID(telegramChatID).ID, price32)
		if err != nil {
			logger.Errorf(err.Error())
		}
	}
}

func getInfoAboutShareByTicker(tracker *tracker.Tracker, ticker string) (*investapi.Share, []string, float32, error) {
	instrResp, err := findShareByTicker(tracker, ticker)
	if instrResp == nil {
		return nil, nil, 0, err
	}

	instrument := instrResp.GetInstrument()
	price, err := m.GetLastPriceByFigi(tracker, instrument)
	if err != nil {
		return nil, nil, 0, err
	}

	shareInfo := printbot.InfoAboutShareByItsTicker(instrument)
	shareInfo += fmt.Sprintf("Нынешняя цена: %.2f ₽\n", price)
	return instrument, []string{shareInfo}, price, nil
}

func findShareByTicker(tracker *tracker.Tracker, ticker string) (*investgo.ShareResponse, error) {
	return tracker.InstrumentsService.ShareByTicker(ticker, "TQBR") //tqbr только для российских компаний
}

func getForecastsAboutShare(tracker *tracker.Tracker, instrument *investapi.Share) []string {
	var forecastResp []string
	forecastResp = append(forecastResp, "Прогнозы от инвест домов:")

	forecast, _ := tracker.InstrumentsService.GetForecastBy(instrument.GetUid())
	for i, target := range forecast.GetTargets() {
		forecastResp = append(forecastResp, printbot.InvestHouseForecast(i, target))
	}
	forecastResp = append(forecastResp, fmt.Sprintf("Согласованный прогноз:%.2f ₽\n", float32(forecast.GetConsensus().GetConsensus().GetUnits())+float32(forecast.GetConsensus().GetConsensus().GetNano())/float32(1000000000)))
	return forecastResp
}
