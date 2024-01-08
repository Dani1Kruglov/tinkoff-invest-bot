package shares

import (
	"fmt"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"tinkoff-investment-bot/internal/database"
	"tinkoff-investment-bot/internal/model"
	printbot "tinkoff-investment-bot/internal/print-bot"
	m "tinkoff-investment-bot/internal/services/marketdata"
)

func ViewInfoOnShareByItsTicker(tracker *model.Tracker, logger *zap.SugaredLogger) {
	instrument, _, err := getInfoAboutShareByTicker(tracker)
	if err != nil {
		logger.Errorf(err.Error())
	}
	getForecastsAboutShare(tracker, instrument)
}

func AddShareToListOfTracked(tracker *model.Tracker, logger *zap.SugaredLogger, db *gorm.DB, telegramID string) {
	instrument, price32, err := getInfoAboutShareByTicker(tracker)
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

		err = database.AddShare(db, &model.Share{
			UID:       instrument.GetUid(),
			Ticker:    instrument.GetTicker(),
			Name:      instrument.GetName(),
			FIGI:      instrument.GetFigi(),
			ClassCode: instrument.GetClassCode(),
		}, database.GetUserByTelegramID(db, telegramID).ID, price32)
		if err != nil {
			logger.Errorf(err.Error())
		}
	}
}

func getInfoAboutShareByTicker(tracker *model.Tracker) (*investapi.Share, float32, error) {
	instrResp, err := findShareByTicker(tracker)
	for err != nil && err.Error() == "rpc error: code = NotFound desc = 50002" {
		fmt.Println("Акции с таким тикером нет, введите другой тикер или '0' чтобы выйти")
		instrResp, err = findShareByTicker(tracker)
	}
	if instrResp == nil {
		return nil, 0, err
	}

	instrument := instrResp.GetInstrument()

	printbot.InfoAboutShareByItsTicker(instrument)

	price, _ := m.GetLastPriceByFigi(tracker, instrument)
	if err != nil {
		return nil, 0, err
	}
	printbot.LastPrice(price)

	fmt.Println("[][][][[][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][]")
	return instrument, price, nil
}

func findShareByTicker(tracker *model.Tracker) (*investgo.ShareResponse, error) {
	ticker, err := printbot.GetTickerFromUser()
	if err != nil {
		return nil, err
	}
	if ticker == "0" {
		return nil, nil
	}
	return tracker.InstrumentsService.ShareByTicker(ticker, "TQBR") //tqbr только для российских компаний
}

func getForecastsAboutShare(tracker *model.Tracker, instrument *investapi.Share) {
	printbot.HeadlineForecastsOfInvestmentHouses()

	forecast, _ := tracker.InstrumentsService.GetForecastBy(instrument.GetUid())
	for i, target := range forecast.GetTargets() {
		printbot.InvestHouseForecast(i, target)
	}
	fmt.Println("======================================")
	printbot.ConsensusForecast(float32(forecast.GetConsensus().GetConsensus().GetUnits()) + float32(forecast.GetConsensus().GetConsensus().GetNano())/float32(1000000000))
}
