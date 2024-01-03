package operations

import (
	"fmt"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	"strconv"
)

func WorkWithOperations(client *investgo.Client, logger *zap.SugaredLogger) {
	operationsService := client.NewOperationsServiceClient()

	portfolioResp, err := operationsService.GetPortfolio(strconv.Itoa(2138000266), pb.PortfolioRequest_RUB)

	if err != nil {
		logger.Errorf(err.Error())
	} else {
		fmt.Printf(" Общая стоимость валют в портфеле. ===== %v\n", portfolioResp.TotalAmountCurrencies.GetUnits())
		//https://russianinvestments.github.io/investAPI/operations/#portfolioresponse
		fmt.Printf("Общая стоимость портфеля. ===== %v\n", portfolioResp.TotalAmountPortfolio.GetUnits())
		fmt.Printf("Массив виртуальных позиций портфеля. ===== %v\n", len(portfolioResp.VirtualPositions))
	}
}
