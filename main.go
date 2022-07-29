package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func getYyyyMmDd(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

func getEnvVariable(key string) string {
	isLocal := os.Getenv("IS_LOCAL")
	if isLocal == "true" {
		fmt.Println("is local")
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("can't load .env file!")
		}
	}

	variable := os.Getenv(key)
	if variable == "" {
		log.Fatal(fmt.Sprintf("can't find %s", key))
	}
	return variable
}

func main() {
	now := time.Now()

	startDate := getYyyyMmDd(now.AddDate(0, 0, -3))
	endDate := getYyyyMmDd(now)

	change := getFluctuationEurJpy(startDate, endDate)

	text := fmt.Sprintf("EUR/JPY rate %s ~ %s\n¥%f => ¥%f\nchange rate: %f%%", startDate, endDate, change.StartRate, change.EndRate, change.ChangePct)
	sendMessageToMoneyChannel(text)

}
