package fiat

import (
	"fmt"
	"time"

	"github.com/mugayoshi/currency_rate_tracker/constants"
	"github.com/mugayoshi/currency_rate_tracker/helpers"
	"github.com/mugayoshi/currency_rate_tracker/slack"
)

func createCurrencyFluctuationNotificationMessage(base string, target string, startDate string, endDate string, data FluctuationDataCurrency) string {
	firstLine := fmt.Sprintf("%s/%s rate", base, target)
	secondLine := fmt.Sprintf("%s ~ %s", startDate, endDate)
	currencySymbol := constants.GetCurrencySymbol(target)
	thirdLine := fmt.Sprintf("%s%f -> %s%f", currencySymbol, data.StartRate, currencySymbol, data.EndRate)
	fourthLine := fmt.Sprintf("change rate: %f%%", data.ChangePct)
	return fmt.Sprintf("%s\n%s\n%s\n%s", firstLine, secondLine, thirdLine, fourthLine)
}

func SendCurrencyUpdates() {
	now := time.Now()

	startDate := helpers.GetYyyyMmDd(now.AddDate(0, 0, -3))
	endDate := helpers.GetYyyyMmDd(now)

	fiat := constants.FIAT_CURRENCY
	eur := fiat.Eur
	jpy := fiat.Jpy
	usd := fiat.Usd
	fmt.Println("EUR/JPY")
	changeEurJpy := getFluctuationBaseJpy(startDate, endDate, eur)
	textEurJpy := createCurrencyFluctuationNotificationMessage(eur, jpy, startDate, endDate, changeEurJpy)
	slack.SendMessageToMoneyChannel(textEurJpy)

	fmt.Println("USD/JPY")
	changeUsdJpy := getFluctuationBaseJpy(startDate, endDate, usd)
	textUsdJpy := createCurrencyFluctuationNotificationMessage(usd, jpy, startDate, endDate, changeUsdJpy)
	slack.SendMessageToMoneyChannel(textUsdJpy)
}
