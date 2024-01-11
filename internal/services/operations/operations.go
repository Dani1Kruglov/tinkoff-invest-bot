package operations

import (
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	"gorm.io/gorm"
	printbot "tinkoff-investment-bot/internal/bot/print"
	"tinkoff-investment-bot/internal/model/tracker"
	"tinkoff-investment-bot/internal/storage"

	u "tinkoff-investment-bot/internal/services/users"
)

func GetUserSecuritiesOnAccount(tracker *tracker.Tracker, logger *zap.SugaredLogger, db *gorm.DB, telegramChatID int64) []string {
	portfolioResp, err := GetPortfolioByAccountID(tracker, db, telegramChatID)
	var userSecurities []string
	if err != nil {
		logger.Errorf(err.Error())
	} else {
		for _, position := range portfolioResp.GetPositions() {
			instrument, err := tracker.InstrumentsService.InstrumentByUid(position.GetInstrumentUid())
			if err != nil {
				logger.Errorf(err.Error())
			}
			userSecurities = append(userSecurities, printbot.InfoAboutUserSecurities(instrument.GetInstrument(), position, portfolioResp.TotalAmountPortfolio.GetUnits()))
		}
	}
	return userSecurities
}

func GetPortfolioByAccountID(tracker *tracker.Tracker, db *gorm.DB, telegramChatID int64) (*investgo.PortfolioResponse, error) {
	accountStorage := storage.NewAccountStorage(db)
	userStorage := storage.NewUserStorage(db)
	account := accountStorage.GetAccountIDByTelegramChatID(telegramChatID)
	if account.AccountID == "" {
		var err error
		account, err = u.GetAccount(tracker)
		if err != nil {
			return nil, err
		}
		err = accountStorage.AddAccount(&account, userStorage.GetUserByTelegramChatID(telegramChatID).ID)
		if err != nil {
			return nil, err
		}
	}
	return tracker.OperationsService.GetPortfolio(account.AccountID, pb.PortfolioRequest_RUB)
}
