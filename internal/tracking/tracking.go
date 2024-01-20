package tracking

import (
	s "tinkoff-investment-bot/internal/model/settings"
	t "tinkoff-investment-bot/internal/model/tracker"
	is "tinkoff-investment-bot/internal/services/instruments/invest-schedules"
	"tinkoff-investment-bot/internal/services/instruments/shares"
	o "tinkoff-investment-bot/internal/services/operations"
)

func TrackByTinkoffToken(settings *s.Settings, tracker *t.Tracker, telegramChatID int64, command string) []string {
	var responses []string

	/*settings.Logger.Infoln("start cron schedule")
	go cron.NewCron(settings.DB, tracker)*/

	switch command {
	case "0":
		break
	case "1":
		responses = o.GetUserSecuritiesOnAccount(tracker, settings.Logger, settings.DB, telegramChatID)
		break
	case "2":
		responses = []string{"Введите тикер акции (MOEX, SBER или другие):"}
		break
	case "3":
		responses = []string{"Введите тикер акции (MOEX, SBER или другие):"}
		break
	case "4":
		responses = is.GetScheduleOnClientSecurities(tracker, settings.Logger, settings.DB, telegramChatID, false)
		break
	case "5":
		responses = is.GetScheduleOnClientSecurities(tracker, settings.Logger, settings.DB, telegramChatID, true)
		break
	default:
		break
	}
	return responses
}

func GetShare(settings *s.Settings, tracker *t.Tracker, ticker string) []string {
	return shares.ViewInfoOnShareByItsTicker(tracker, settings.Logger, ticker)
}
