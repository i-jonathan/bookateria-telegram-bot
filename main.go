package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/yoruba-codigy/goTelegram"
	"io"
	"log"
	"net/http"
)

type setWebHook struct {
	Url string `json:"url"`
}

type getWebHookInfo struct {
	Ok     bool `json:"ok"`
	Result Result `json:"result"`
}

type Result struct {
	Url                  string `json:"url"`
	HasCustomCertificate bool   `json:"has_custom_certificate"`
	PendingUpdateCount   int    `json:"pending_update_count"`
	MaxConnections       int    `json:"max_connections"`
	IpAddress            string `json:"ip_address"`
}

var bot goTelegram.Bot
var apiURL = "https://bookateria-api.herokuapp.com/v1/"

func main() {
	set := setWebhook("https://c0eeb7bdaded.ngrok.io")
	fmt.Println(set)
	var err error
	bot, err = goTelegram.NewBot("891332272:AAHJ8fm_iwdHNO1VhNb0g2BVldpQx35Ooms")

	if err != nil {
		log.Println(err)
	}
	bot.SetHandler(handler)
	log.Println("Starting Server")
	err = http.ListenAndServe(":5000", http.HandlerFunc(bot.UpdateHandler))

	if err != nil {
		log.Println("Failed")
		log.Fatalln(err)
		return
	}
}

func setWebhook(webHookURL string) bool {
	url := "https://api.telegram.org/bot891332272:AAHJ8fm_iwdHNO1VhNb0g2BVldpQx35Ooms/"

	resp, err := http.Get(url + "getWebhookInfo")

	if err != nil {
		log.Println(err)
		return false
	}

	var webhookInfo = &getWebHookInfo{}

	if err := json.NewDecoder(resp.Body).Decode(webhookInfo); err != nil {
		log.Println(err)
		return false
	}

	log.Println(webhookInfo.Result.Url)

	if webHookURL == webhookInfo.Result.Url {
		return true
	}

	body := setWebHook{Url: webHookURL}

	jsonBody, err := json.Marshal(body)

	if err != nil {
		log.Println(err)
		return false
	}

	resp, err = http.Post(url + "setWebhook", "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		log.Println(err)
		return false
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Println(resp.StatusCode)
		log.Println("Status not okay")
		return false
	}

	return true
}