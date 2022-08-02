package fiat

import (
	"fmt"
	"time"

	"github.com/mugayoshi/currency_rate_tracker/slack"
)

type Currency string

const (
	EUR Currency = "EUR"
	JPY Currency = "JPY"
	USD Currency = "USD"
)

func getYyyyMmDd(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}
func createCurrencyFluctuationNotificationMessage(base Currency, target Currency, startDate string, endDate string, data FluctuationDataCurrency) string {
	firstLine := fmt.Sprintf("%s/%s rate", base, target)
	secondLine := fmt.Sprintf("%s ~ %s", startDate, endDate)
	getCurrencySymbol := func(currency Currency) string {
		switch currency {
		case "EUR":
			return "€"
		case "JPY":
			return "¥"
		case "USD":
			return "$"
		default:
			return ""
		}
	}
	currencySymbol := getCurrencySymbol(target)
	thirdLine := fmt.Sprintf("%s%f -> %s%f", currencySymbol, data.StartRate, currencySymbol, data.EndRate)
	fourthLine := fmt.Sprintf("change rate: %f%%", data.ChangePct)
	return fmt.Sprintf("%s\n%s\n%s\n%s", firstLine, secondLine, thirdLine, fourthLine)
}

func SendCurrencyUpdates() {
	now := time.Now()

	startDate := getYyyyMmDd(now.AddDate(0, 0, -3))
	endDate := getYyyyMmDd(now)

	fmt.Println("EUR/JPY")
	changeEurJpy := getFluctuationBaseJpy(startDate, endDate, EUR)
	textEurJpy := createCurrencyFluctuationNotificationMessage(EUR, JPY, startDate, endDate, changeEurJpy)
	slack.SendMessageToMoneyChannel(textEurJpy)

	fmt.Println("USD/JPY")
	changeUsdJpy := getFluctuationBaseJpy(startDate, endDate, USD)
	textUsdJpy := createCurrencyFluctuationNotificationMessage(USD, JPY, startDate, endDate, changeUsdJpy)
	slack.SendMessageToMoneyChannel(textUsdJpy)
}
