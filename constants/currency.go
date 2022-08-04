package constants

type FiatCurrency struct {
	Eur string
	Jpy string
	Usd string
}

type CryptoCurrency struct {
	Btc string
	Eth string
}

var (
	CRYPTO_CURRENCY = CryptoCurrency{Btc: "BTC", Eth: "ETH"}
	FIAT_CURRENCY   = FiatCurrency{Eur: "EUR", Jpy: "JPY", Usd: "USD"}
)

func GetCurrencySymbol(currency string) string {
	fiat := FIAT_CURRENCY
	switch currency {
	case fiat.Eur:
		return "€"
	case fiat.Jpy:
		return "¥"
	case fiat.Usd:
		return "$"
	default:
		return ""
	}
}
