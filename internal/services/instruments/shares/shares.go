package shares

import (
	"fmt"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	"strconv"
	printbot "tinkoff-investment-bot/internal/bot/print"
	"tinkoff-investment-bot/internal/model/database"
	"tinkoff-investment-bot/internal/model/settings"
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

func GetShareForFavoriteList(tracker *tracker.Tracker, settings *settings.Settings, ticker string) []string {
	_, shareResp, _, err := getInfoAboutShareByTicker(tracker, ticker)
	if err != nil {
		settings.Logger.Errorf(err.Error())
	}
	shareResp[0] += "\nДобавить акцию в список отслеживаемых?\n"
	return shareResp
}

func AddShareToListOfTracked(tracker *tracker.Tracker, settings *settings.Settings, ticker string, telegramChatID int64, customPrice string) []string {
	instrument, _, price32, err := getInfoAboutShareByTicker(tracker, ticker)
	if err != nil {
		settings.Logger.Errorf(err.Error())
	}

	if customPrice != "" {
		price, err := strconv.ParseFloat(customPrice, 32)
		if err != nil {
			settings.Logger.Errorf(err.Error())
		}
		price32 = float32(price)
	}
	userStorage := storage.NewUserStorage(settings.DB)
	shareStorage := storage.NewShareStorage(settings.DB)

	err = shareStorage.AddShare(&database.Share{
		UID:       instrument.GetUid(),
		Ticker:    instrument.GetTicker(),
		Name:      instrument.GetName(),
		FIGI:      instrument.GetFigi(),
		ClassCode: instrument.GetClassCode(),
	}, userStorage.GetUserByTelegramChatID(telegramChatID).ID, price32)
	if err != nil {
		settings.Logger.Errorf(err.Error())
	}

	return []string{"Акция успешно добавлена в отслеживаемый список"}
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
