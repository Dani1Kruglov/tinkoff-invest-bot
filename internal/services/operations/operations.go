package operations

import (
	"fmt"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"

	"tinkoff-investment-bot/internal/model"
	u "tinkoff-investment-bot/internal/services/users"
)

func GetUserSharesOnAccount(logger *zap.SugaredLogger, tracker *model.Tracker) {
	portfolioResp, err := GetPortfolioByAccountID(tracker)

	if err != nil {
		logger.Errorf(err.Error())
	} else {
		for _, position := range portfolioResp.GetPositions() {
			fmt.Println("-  -  -  -  -  -  -  -  -  -  -")
			instrument, err := tracker.InstrumentsService.FindInstrument(position.GetInstrumentUid())
			if err != nil {
				logger.Errorf(err.Error())
			}
			fmt.Printf("Название ценной бумаги: %v\n", instrument.GetInstruments()[0].GetName())
			fmt.Printf("Тикер ценной бумаги: %v\n", instrument.GetInstruments()[0].GetTicker())
			fmt.Printf("Тип инструмента: %v\n", position.GetInstrumentType())
			fmt.Printf("Количество в портфеле в штуках: %d.%d\n", position.GetQuantity().GetUnits(), position.GetQuantity().GetNano()/10000000)
			fmt.Printf("Цена за шт в портфеле: %d.%d ₽\n", position.GetAveragePositionPrice().GetUnits(), position.GetAveragePositionPriceFifo().GetNano()/10000000)
			fmt.Println("...........Новое сообщение.........")
			fmt.Printf("Цена ценной бумаги сейчас: %d.%d ₽\n", position.GetCurrentPrice().GetUnits(), position.GetCurrentPrice().GetNano()/10000000)
		}
		fmt.Println("////////////////////////////////////")
		fmt.Printf("Общая стоимость портфеля: %v ₽\n", portfolioResp.TotalAmountPortfolio.GetUnits())
	}
	fmt.Println("[][][][[][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][]")
}

func GetPortfolioByAccountID(tracker *model.Tracker) (*investgo.PortfolioResponse, error) {
	accountID, err := u.GetAccountID(tracker)
	if err != nil {
		return nil, err
	}
	return tracker.OperationsService.GetPortfolio(accountID, pb.PortfolioRequest_RUB)
}
