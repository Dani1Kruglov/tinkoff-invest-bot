package operations

import (
	"fmt"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	"gorm.io/gorm"
	printbot "tinkoff-investment-bot/internal/print-bot"
	"tinkoff-investment-bot/internal/storage"

	"tinkoff-investment-bot/internal/model"
	u "tinkoff-investment-bot/internal/services/users"
)

func GetUserSecuritiesOnAccount(tracker *model.Tracker, logger *zap.SugaredLogger, db *gorm.DB, telegramID string) {
	portfolioResp, err := GetPortfolioByAccountID(tracker, db, telegramID)
	if err != nil {
		logger.Errorf(err.Error())
	} else {
		for _, position := range portfolioResp.GetPositions() {
			instrument, err := tracker.InstrumentsService.InstrumentByUid(position.GetInstrumentUid())
			if err != nil {
				logger.Errorf(err.Error())
			}
			printbot.InfoAboutUserSecurities(instrument.GetInstrument(), position)
		}
		fmt.Println("////////////////////////////////////")
	}
	printbot.TotalAmountPortfolioUser(portfolioResp)
	fmt.Println("[][][][[][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][]")
}

func GetPortfolioByAccountID(tracker *model.Tracker, db *gorm.DB, telegramID string) (*investgo.PortfolioResponse, error) {
	accountStorage := storage.NewAccountStorage(db)
	userStorage := storage.NewUserStorage(db)
	account := accountStorage.GetAccountIDByTelegramID(telegramID)
	if account.AccountID == "" {
		var err error
		account, err = u.GetAccount(tracker)
		if err != nil {
			return nil, err
		}
		err = accountStorage.AddAccount(&account, userStorage.GetUserByTelegramID(telegramID).ID)
		if err != nil {
			return nil, err
		}
	}
	return tracker.OperationsService.GetPortfolio(account.AccountID, pb.PortfolioRequest_RUB)
}
