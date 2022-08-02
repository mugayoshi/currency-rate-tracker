package cryptoCurrency

import (
	"fmt"
	"time"

	"github.com/mugayoshi/currency_rate_tracker/slack"
	"golang.org/x/text/message"
)

func createNotificationMessage(timestamp int64, btcPrice float32, ethPrice float32) string {
	now := time.Unix(timestamp, 0).UTC()
	formattedDate := now.Format("2006 2 Jan Mon")
	englishPrinter := message.NewPrinter(message.MatchLanguage("en"))
	btc := englishPrinter.Sprintf("¥%.3f", btcPrice)
	eth := englishPrinter.Sprintf("¥%.3f", ethPrice)
	return fmt.Sprintf("%s\nBTC: %s\nETH: %s", formattedDate, btc, eth)
}

func CheckRates() {
	fmt.Println("check crypto currency rate")
	symbols := []string{"BTC", "ETH"}
	data := getLiveData("JPY", symbols)

	text := createNotificationMessage(data.Timestamp, data.Rates.Btc, data.Rates.Eth)
	slack.SendMessageToMoneyChannel(text)
}
