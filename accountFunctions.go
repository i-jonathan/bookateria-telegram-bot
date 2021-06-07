package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/yoruba-codigy/goTelegram"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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
		StayIn:   true,
	}

	postBody, err := json.Marshal(loginDetails)

	if err != nil {
		log.Println(err)
	}

	var resp *http.Response

	resp, err = http.Post(apiURL+"auth/login", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		log.Println(err)
	}

	var loginData LoginResponse

	if err = json.NewDecoder(resp.Body).Decode(&loginData); err != nil {
		log.Println(err)
	}

	log.Println(loginData.Value)

	// Save token to redis
	var ctx = context.Background()

	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB: 0,
	})

	redisTime := 720 * time.Hour

	err = client.Set(ctx, strconv.Itoa(update.UserID), loginData.Value, redisTime).Err()

	if err != nil {
		bot.EditMessage(replyUpdate.Message, "Login Unsuccessful. Try again shortly")
	}
	bot.EditMessage(replyUpdate.Message, "Signed in Successfully.")
}
