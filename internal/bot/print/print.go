package print

import (
	"fmt"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"time"
)

func InfoAboutUserSecurities(instrument *investapi.Instrument, position *investapi.PortfolioPosition, totalAmountPortfolio int64) string {
	result := ""
	result += fmt.Sprintf("Название ценной бумаги: %v\n", instrument.GetName())
	result += fmt.Sprintf("Тикер ценной бумаги: %v\n", instrument.GetTicker())
	result += fmt.Sprintf("Тип инструмента: %v\n", position.GetInstrumentType())
	result += fmt.Sprintf("Количество в портфеле в штуках: %d.%d\n", position.GetQuantity().GetUnits(), position.GetQuantity().GetNano()/10000000)
	result += fmt.Sprintf("Цена за шт в портфеле: %d.%d ₽\n", position.GetAveragePositionPrice().GetUnits(), position.GetAveragePositionPriceFifo().GetNano()/10000000)
	result += fmt.Sprintf("Цена бумаги сейчас: %d.%d ₽\n", position.GetCurrentPrice().GetUnits(), position.GetCurrentPrice().GetNano()/10000000)
	result += fmt.Sprintf("Общая стоимость портфеля: %v ₽\n", totalAmountPortfolio)
	return result
}

func InfoAboutShareByItsTicker(instrument *investapi.Share) string {
	result := ""
	result += fmt.Sprintf("Название акции: %v\n", instrument.GetName())
	result += fmt.Sprintf("Тикер акции: %v\n", instrument.GetTicker())
	result += fmt.Sprintf("Класс код акции: %v\n", instrument.GetClassCode())
	return result
}

func InvestHouseForecast(i int, target *investapi.GetForecastResponse_TargetItem) string {
	result := ""
	result += fmt.Sprintf("Номер компании: %d\n", i+1)
	result += fmt.Sprintf("Название компании, давшей прогноз: %v\n", target.GetCompany())
	result += fmt.Sprintf("Прогноз: %v\n", target.GetRecommendation())
	result += fmt.Sprintf("Прогнозируемая цена: %v.%v ₽\n", target.GetTargetPrice().GetUnits(), target.GetTargetPrice().GetNano()/10000000)
	result += fmt.Sprintf("Изменение цены: %v.%v ₽\n", target.GetPriceChange().GetUnits(), target.GetPriceChange().GetNano()/10000000)
	return result
}

func UserAccountSelect(accounts []*investapi.Account) (int, error) {
	var num int
	fmt.Println("Выберите счёт, который хотите просмотреть: ")
	for i, account := range accounts {
		fmt.Println("-  -  -  -  -  -  -  -  -  -  -")
		fmt.Printf("%d: \n id счёта: %v\n", i+1, account.GetId())
		fmt.Printf("id счёта: %v\n", account.GetName())
	}
	_, err := fmt.Scan(&num)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func PriceChange(currentPrice float32, pastPrice float32, message string) {
	fmt.Println("-  -  -  -  -  -  -  -  -  -  -")
	fmt.Println(time.Now())
	fmt.Println(message)
	fmt.Printf("Прошлая цена: %.2f\n", pastPrice)
	fmt.Printf("Нынешняя цена: %.2f\n", currentPrice)
}
