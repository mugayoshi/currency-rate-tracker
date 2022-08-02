package cryptoCurrency

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mugayoshi/currency_rate_tracker/helpers"
)

type CryptoCurrencyRates struct {
	Btc float32 `json:"BTC"`
	Eth float32 `json:"ETH"`
}

type CryptoCurrencyLiveData struct {
	Target    string              `json:"target"`
	Rates     CryptoCurrencyRates `json:"rates"`
	Timestamp int64               `json:"timestamp"`
}

func getLiveData(target string, symbols []string) CryptoCurrencyLiveData {
	symbolArray := strings.Join(symbols, ",")
	accessKey := helpers.GetEnvVariable("COINLAYER_API_ACCESS_KEY")
	apiPath := "live"
	apiParams := fmt.Sprintf("access_key=%s&target=%s&symbols=%s", accessKey, target, symbolArray)
	body := callCoinlayerApi(apiPath, apiParams)
	var result CryptoCurrencyLiveData

	json.Unmarshal(body, &result)
	return result
}

func callCoinlayerApi(path string, params string) []byte {
	url := fmt.Sprintf("http://api.coinlayer.com//%s?%s", path, params)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return body
}
