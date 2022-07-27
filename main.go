package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func getYyyyMmDd(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("can't load .env file!")
	}
	return os.Getenv(key)
}

func main() {
	now := time.Now()

	startDate := getYyyyMmDd(now)
	endDate := getYyyyMmDd(now.AddDate(0, 0, -3))

	url := fmt.Sprintf("https://api.apilayer.com/fixer/fluctuation?start_date=%s&end_date=%s&base=EUR&symbols=JPY,USD", startDate, endDate)
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

	fmt.Println(string(body))
}
