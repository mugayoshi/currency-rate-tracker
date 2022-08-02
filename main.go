package main

import (
	"flag"

	cryptoCurrency "github.com/mugayoshi/currency_rate_tracker/currency/crypto"
	"github.com/mugayoshi/currency_rate_tracker/currency/fiat"
)

func main() {

	f := flag.String("type", "", "command type")
	flag.Parse()
	commandType := *f
	switch commandType {
	case "fiat":
		fiat.SendCurrencyUpdates()
		return
	case "crypto":
		cryptoCurrency.CheckRates()
		return

	}

}
