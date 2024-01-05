package print_bot

import (
	"fmt"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"time"
)

func GetTokenFromUser() (string, error) {
	fmt.Println("Введите 'токен для чтения' своего брокерского счета:")
	var token string
	_, err := fmt.Scan(&token)
	if err != nil {
		_ = fmt.Errorf(err.Error())
	}
	return token, nil
}

func MainMenu() string {
	var command string
	fmt.Println("Введите что хотите сделать: ")
	fmt.Println("0 - Закончить")
	fmt.Println("1 - Посмотреть свои ценные бумаги")
	fmt.Println("2 - Посмотреть данные по акции по ее тикеру")
	fmt.Println("3 - Добавить акцию не из вашего портфеля в список отслеживаемых ")
	fmt.Println("4 - Посмотреть расписание дивидендов по вашим ценным бумагам")
	fmt.Println("5 - Посмотреть расписание отчётов по вашим ценным бумагам")
	_, err := fmt.Scan(&command)
	if err != nil {
		_ = fmt.Errorf(err.Error())
	}
	return command
}

func InfoAboutUserSecurities(instrument *investapi.Instrument, position *investapi.PortfolioPosition) { // цикл не в этой функции потому что в тг будет каждая бумага отдельным сообщением
	fmt.Println("-  -  -  -  -  -  -  -  -  -  -")
	fmt.Printf("Название ценной бумаги: %v\n", instrument.GetName())
	fmt.Printf("Тикер ценной бумаги: %v\n", instrument.GetTicker())
	fmt.Printf("Тип инструмента: %v\n", position.GetInstrumentType())
	fmt.Printf("Количество в портфеле в штуках: %d.%d\n", position.GetQuantity().GetUnits(), position.GetQuantity().GetNano()/10000000)
	fmt.Printf("Цена за шт в портфеле: %d.%d ₽\n", position.GetAveragePositionPrice().GetUnits(), position.GetAveragePositionPriceFifo().GetNano()/10000000)
	fmt.Println("...........Новое сообщение.........")
	fmt.Printf("Цена ценной бумаги сейчас: %d.%d ₽\n", position.GetCurrentPrice().GetUnits(), position.GetCurrentPrice().GetNano()/10000000)
}

func TotalAmountPortfolioUser(portfolioResp *investgo.PortfolioResponse) {
	fmt.Printf("Общая стоимость портфеля: %v ₽\n", portfolioResp.TotalAmountPortfolio.GetUnits())
}

func GetTickerFromUser() (string, error) {
	fmt.Println("Введите тикер акции (MOEX, SBER и так далее)")
	var ticker string
	_, err := fmt.Scan(&ticker)
	if err != nil {
		return "", err
	}
	return ticker, nil
}

func InfoAboutShareByItsTicker(instrument *investapi.Share) {
	fmt.Println("-  -  -  -  -  -  -  -  -  -  -")
	fmt.Printf("Название акции: %v\n", instrument.GetName())
	fmt.Printf("Тикер акции: %v\n", instrument.GetTicker())
	fmt.Printf("Класс код акции: %v\n", instrument.GetClassCode())
}

func LastPrice(units int64, nano int32) {
	fmt.Printf("Нынешняя цена: %d.%d ₽\n", units, nano)
}

func HeadlineForecastsOfInvestmentHouses() {
	fmt.Println("Прогнозы от инвест домов:")
}

func InvestHouseForecast(i int, target *investapi.GetForecastResponse_TargetItem) {
	fmt.Println("-  -  -  -  -  -  -  -  -  -  -")
	fmt.Printf("Номер компании: %d\n", i+1)
	fmt.Printf("Название компании, давшей прогноз: %v\n", target.GetCompany())
	fmt.Printf("Прогноз: %v\n", target.GetRecommendation())
	fmt.Printf("Прогнозируемая цена: %v.%v ₽\n", target.GetTargetPrice().GetUnits(), target.GetTargetPrice().GetNano()/10000000)
	fmt.Printf("Изменение цены: %v.%v ₽\n", target.GetPriceChange().GetUnits(), target.GetPriceChange().GetNano()/10000000)
}

func ConsensusForecast(units int64, nano int32) {
	fmt.Printf("Согласованный прогноз: %v.%v ₽\n", units, nano)
}

func PrintScheduleOfReports(name string, reports *investgo.GetAssetReportsResponse) {
	fmt.Printf("Отчет по бумаге %v в ближайшие 6 месяцев: %v\n", name, reports.GetAssetReportsResponse)
}
func PrintScheduleOfDividend(name string, timeLocal time.Time, dividendNet *investapi.MoneyValue) {
	fmt.Printf("Дивиденды по бумаге '%v' в ближайшие 6 месяцев:\nДата фактических выплат по МСК: %v\nВеличина дивиденда на 1 ценную бумагу (включая валюту): %d.%d\n",
		name, timeLocal, dividendNet.GetUnits(), dividendNet.GetNano())
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
