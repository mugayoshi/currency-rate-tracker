package fiat

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/mugayoshi/currency_rate_tracker/aws"
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

func sendUpdateToSlack(startDate, endDate, targetCurrency, baseCurrency string, dbClient *dynamodb.Client) {
	change := getFluctuationBaseJpy(startDate, endDate, targetCurrency)
	text := createCurrencyFluctuationNotificationMessage(targetCurrency, baseCurrency, startDate, endDate, change)
	threshold := func(target string) float32 {
		t := helpers.GetEnvVariable(fmt.Sprintf("THRESHOLD_%s", targetCurrency))
		f, e := strconv.ParseFloat(t, 32)
		if e != nil {
			fmt.Println("cannot get threshold value")
			return 0
		}
		return float32(f)
	}(targetCurrency)

	if change.EndRate < threshold {
		text = fmt.Sprintf("%s\nit's below %f", text, threshold)
	}
	currentItem, getErr := aws.GetCurrencyItem(dbClient, targetCurrency)
	if getErr != nil {
		fmt.Println("couldn't get the last data")
		return
	}

	isMinimumRate := currentItem.Minimum.Rate > change.EndRate
	isSuccess, updateErr := aws.UpdateLastData(dbClient, targetCurrency, change.EndRate, endDate)
	if updateErr != nil || !isSuccess {
		fmt.Println("couldn't update the last data")
		return
	}
	if isMinimumRate {
		text = fmt.Sprintf("%s\n!!Minimum rate: %f!!", text, change.EndRate)
		aws.UpdateMinimumRate(dbClient, targetCurrency, change.EndRate, endDate)
	}
	slack.SendMessageToMoneyChannel(text)

}

func SendCurrencyUpdates() {
	now := time.Now()

	startDate := helpers.GetYyyyMmDd(now.AddDate(0, 0, -3))
	endDate := helpers.GetYyyyMmDd(now)

	fiat := constants.FIAT_CURRENCY
	jpy := fiat.Jpy
	dbClient := aws.GetDynamoDbClient()
	targetCurrencies := []string{fiat.Eur, fiat.Usd}

	var wg sync.WaitGroup
	for _, c := range targetCurrencies {
		wg.Add(1)
		go func(target string) {
			defer wg.Done()
			sendUpdateToSlack(startDate, endDate, target, jpy, dbClient)
		}(c)
	}
	wg.Wait()
}
