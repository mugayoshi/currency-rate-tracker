package cryptoCurrency

import (
	"fmt"
	"time"

	"github.com/mugayoshi/currency_rate_tracker/constants"
	"github.com/mugayoshi/currency_rate_tracker/helpers"
	"github.com/mugayoshi/currency_rate_tracker/slack"
	"golang.org/x/text/message"
)

func createNotificationMessage(timestamp int64, btcPrice float32, ethPrice float32) string {
	now := time.Unix(timestamp, 0).UTC()
	yyyyDdMm := now.Format("2006 2 Jan Mon")
	hhMmSs := now.Format("15:04:05")
	englishPrinter := message.NewPrinter(message.MatchLanguage("en"))
	btc := englishPrinter.Sprintf("¥%.3f", btcPrice)
	eth := englishPrinter.Sprintf("¥%.3f", ethPrice)
	return fmt.Sprintf("%s %s\nBTC: %s\nETH: %s", yyyyDdMm, hhMmSs, btc, eth)
}

type PriceDiffDataForNotification struct {
	fiat       string
	crypto     string
	startDate  string
	endDate    string
	startPrice float32
	endPrice   float32
}

func createNotificationMessageTrend(data PriceDiffDataForNotification) string {
	currencySymbol := constants.GetCurrencySymbol(data.fiat)
	startPriceStr := helpers.CommaSplitString(data.startPrice)
	endPriceStr := helpers.CommaSplitString(data.endPrice)
	firstLine := fmt.Sprintf("%s price trend", data.crypto)
	secondLine := fmt.Sprintf("%s~%s", data.startDate, data.endDate)
	thirdLine := fmt.Sprintf("%s%s->%s%s", currencySymbol, startPriceStr, currencySymbol, endPriceStr)
	return fmt.Sprintf("%s\n%s\n%s", firstLine, secondLine, thirdLine)
}

func CheckCurrentRates() {
	fmt.Println("check crypto currency rate")
	crypto := constants.CRYPTO_CURRENCY
	symbols := []string{crypto.Btc, crypto.Eth}
	data := getLiveData(constants.FIAT_CURRENCY.Jpy, symbols)

	text := createNotificationMessage(data.Timestamp, data.Rates.Btc, data.Rates.Eth)
	slack.SendMessageToMoneyChannel(text)
}

// compare the prices of BTC and ETH between now and 3 days ago
func CheckCurrentAndPreviousPrices(days int) {
	fmt.Printf("check prices between now and %d days ago", days)
	now := time.Now()
	target := constants.FIAT_CURRENCY.Jpy
	crypto := constants.CRYPTO_CURRENCY
	symbols := []string{crypto.Btc, crypto.Eth}
	startDate := helpers.GetYyyyMmDd(now.AddDate(0, 0, days*-1))
	endDate := helpers.GetYyyyMmDd(now)

	dataNow := getHistoricalData(endDate, target, symbols)
	dataDaysAgo := getHistoricalData(startDate, target, symbols)
	var btcData = PriceDiffDataForNotification{fiat: target, crypto: crypto.Btc, startDate: startDate, endDate: endDate, startPrice: dataDaysAgo.Rates.Btc, endPrice: dataNow.Rates.Btc}
	notificationBtc := createNotificationMessageTrend(btcData)
	slack.SendMessageToMoneyChannel(notificationBtc)

	var ethData = PriceDiffDataForNotification{fiat: target, crypto: crypto.Eth, startDate: startDate, endDate: endDate, startPrice: dataDaysAgo.Rates.Eth, endPrice: dataNow.Rates.Eth}
	notificationEth := createNotificationMessageTrend(ethData)
	slack.SendMessageToMoneyChannel(notificationEth)
}
