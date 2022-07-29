package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func sendMessageToMoneyChannel(message string) {
	messageValue := fmt.Sprintf(`{"text":"%s"}`, message)
	var jsonStr = []byte(messageValue)
	url := getEnvVariable("SLACK_WEBHOOK_MONEY")
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
	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
