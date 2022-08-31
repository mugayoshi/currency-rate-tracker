package slack

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/mugayoshi/currency_rate_tracker/helpers"
)

func SendMessageToMoneyChannel(message string) {
	messageValue := fmt.Sprintf(`{"text":"%s"}`, message)
	var jsonStr = []byte(messageValue)
	url := helpers.GetEnvVariable("SLACK_WEBHOOK_MONEY")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	log.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))
	if strings.Contains(resp.Status, "200") {
		log.Println("successfully sent a message to Slack.")
	}
}
