package operations

import (
	"fmt"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	printbot "tinkoff-investment-bot/internal/print-bot"

	"tinkoff-investment-bot/internal/model"
	u "tinkoff-investment-bot/internal/services/users"
)

func GetUserSecuritiesOnAccount(tracker *model.Tracker, logger *zap.SugaredLogger) {
	portfolioResp, err := GetPortfolioByAccountID(tracker)
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

func GetPortfolioByAccountID(tracker *model.Tracker) (*investgo.PortfolioResponse, error) {
	accountID, err := u.GetAccountID(tracker)
	if err != nil {
		return nil, err
	}
	return tracker.OperationsService.GetPortfolio(accountID, pb.PortfolioRequest_RUB)
}
