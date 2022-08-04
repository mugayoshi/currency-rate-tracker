package main

import (
	"flag"
	"fmt"

	cryptoCurrency "github.com/mugayoshi/currency_rate_tracker/currency/crypto"
	"github.com/mugayoshi/currency_rate_tracker/currency/fiat"
)

type CommandArgs struct {
	currencyType string
	frequency    string
}
type CurrencyTypeArg struct {
	fiat   string
	crypto string
}

type FrequencyArg struct {
	high    string
	middle  string
	low     string
	daily   string
	weekly  string
	monthly string
}

var (
	COMMAND_ARGS  = CommandArgs{currencyType: "type", frequency: "frequency"}
	COMMAND_TYPE  = CurrencyTypeArg{fiat: "fiat", crypto: "crypto"}
	FREQUENCY_ARG = FrequencyArg{high: "high", middle: "middle", low: "low", daily: "daily", weekly: "weekly", monthly: "monthly"}
)

func runCommandCrypto(frequency string) {
	switch frequency {
	case FREQUENCY_ARG.high:
		cryptoCurrency.CheckCurrentRates()
		return
	case FREQUENCY_ARG.weekly:
		cryptoCurrency.CheckCurrentAndPreviousPrices(7)
		return
	default:
		fmt.Println(frequency)
		return
	}
}

func runCommandFiat(frequency string) {
	switch frequency {
	case FREQUENCY_ARG.high:
		fiat.SendCurrencyUpdates()
		return
	default:
		fmt.Println(frequency)
		return
	}
}

func main() {

	ct := flag.String(COMMAND_ARGS.currencyType, "", "command type")
	f := flag.String(COMMAND_ARGS.frequency, "", "frequency")
	flag.Parse()
	// commandType := *ct
	switch *ct {
	case COMMAND_TYPE.fiat:
		runCommandFiat(*f)
		return
	case COMMAND_TYPE.crypto:
		runCommandCrypto(*f)
	}

}
