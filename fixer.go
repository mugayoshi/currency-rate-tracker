package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FluctuationDataCurrency struct {
	StartRate float32 `json:"start_rate"`
	EndRate   float32 `json:"end_rate"`
	Change    float32 `json:"change"`
	ChangePct float32 `json:"change_pct"`
}

type FluctuationRate struct {
	Usd FluctuationDataCurrency `json:"USD,omitempty"`
	Jpy FluctuationDataCurrency `json:"JPY,omitempty"`
}

type FluctuationResponse struct {
	Success     bool            `json:"success"`
	Fluctuation string          `json:"fluctuation"`
	StartDate   string          `json:"start_date"`
	EndDate     string          `json:"end_date"`
	Base        string          `json:"base"`
	Rates       FluctuationRate `json:"rates"`
}

func getFluctuationEurJpy(startDate string, endDate string) FluctuationDataCurrency {

	params := fmt.Sprintf("start_date=%s&end_date=%s&base=EUR&symbols=JPY", startDate, endDate)
	body := callFixerApi("fluctuation", params)
	fmt.Println(string(body))
	var result FluctuationResponse
	json.Unmarshal(body, &result)

	return result.Rates.Jpy
}

func callFixerApi(path string, params string) []byte {
	url := fmt.Sprintf("https://api.apilayer.com/fixer/%s?%s", path, params)
	fmt.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	apiKey := getEnvVariable("FIXER_API_KEY")
	req.Header.Set("apikey", apiKey)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return body
}
