package main

import (
	"bytes"
	"encoding/json"
	"github.com/yoruba-codigy/goTelegram"
	"log"
	"net/http"
	"strings"
)

func handleLogin(update query) {
	//Do Random Login Stuff Here

	defer func() {
		replies = remove(replies, update)
	}()

	var replyUpdate goTelegram.Update

	replyUpdate.Message.MessageID = update.MessageID
	replyUpdate.Message.Chat.ID = update.ChatID

	// Get username and password from message text
	messageParts := strings.Split(update.Text, "\n")

	loginDetails := Login{
		Email:    messageParts[0],
		Password: messageParts[1],
		//StayIn:   true,
	}

	postBody, err := json.Marshal(loginDetails)

	if err != nil {
		log.Println(err)
	}

	var resp *http.Response

	resp, err = http.Post(apiURL + "auth/login", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		log.Println(err)
	}

	var loginData LoginResponse

	if err = json.NewDecoder(resp.Body).Decode(&loginData); err != nil {
		log.Println(err)
	}

	log.Println(loginData.Value)
	// Save token in redis and all of that
	// For now, return the token

	bot.EditMessage(replyUpdate.Message, loginData.Value)
}
