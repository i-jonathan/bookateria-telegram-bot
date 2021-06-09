package main

import (
	"github.com/yoruba-codigy/goTelegram"
	"log"

	//"net/http"
	"strings"
)

func handler(update goTelegram.Update) {

	switch update.Type {
	case "text":
		go processRequest(update)
	case "callback":
		go processCallback(update)
	}
}

func processRequest(update goTelegram.Update) {

	if reply, ok := get(replies, update.Message.Chat.ID); ok {
		if reply.UserID == update.Message.From.ID {
			reply.Text = update.Message.Text
			switch reply.Type {
			case "search":
				search(reply)
			case "login":
				handleLogin(reply)
			}
		}
	}

	chat := update.Message.Chat

	parts := strings.Fields(update.Message.Text)
	var command string

	if len(parts) > 0 {
		if update.Message.Chat.Type != "private" && strings.HasSuffix(parts[0], "bookateria_bot") {
			command = strings.Split(parts[0], "@")[0]
		} else if update.Message.Chat.Type == "private" {
			command = parts[0]
		}
	} else {
		return
	}

	switch command {
	case "/start":
		//bot.SendMessage("Hello", chat)
		bot.DeleteKeyboard()

		if update.Message.Chat.Type == "private" {
			bot.AddButton("Account", "account")
		}

		bot.AddButton("Documents", "documents")
		bot.MakeKeyboard(2)
		err := bot.SendMessage("Hello. Where would you like to go?", chat)
		if err != nil {
			log.Println(err)
		}
	}
}

func processCallback(update goTelegram.Update) {
	defer func(bot *goTelegram.Bot, callbackID string) {
		err := bot.AnswerCallback(callbackID)
		if err != nil {
			log.Println(err)
		}
	}(&bot, update.CallbackQuery.ID)

	var command string

	if strings.HasPrefix(update.CallbackQuery.Data, "docID") {
		command = "docID"
	} else if strings.HasPrefix(update.CallbackQuery.Data, "all") {
		command = "all"
	} else if strings.HasPrefix(update.CallbackQuery.Data, "contsearch") {
		command = "contsearch"
	} else {
		command = update.CallbackQuery.Data
	}

	switch command {
	case "account":
		bot.DeleteKeyboard()
		bot.AddButton("Login", "login")
		bot.AddButton("Register", "register")
		bot.AddButton("Logout", "logout")
		bot.MakeKeyboard(3)
		err := bot.EditMessage(update.CallbackQuery.Message, "Accounts")
		if err != nil {
			log.Println(err)
		}

	case "documents":
		bot.DeleteKeyboard()

		if update.CallbackQuery.Message.Chat.Type == "private" {
			bot.AddButton("Add", "add")
			bot.AddButton("Update", "update")
			bot.AddButton("My Uploads", "mine")
		}

		bot.AddButton("All", "all-1")
		bot.AddButton("Search", "search")
		bot.AddButton("Tags", "tags")
		bot.AddButton("Categories", "cat")

		bot.MakeKeyboard(3)
		err := bot.EditMessage(update.CallbackQuery.Message, "Documents")
		if err != nil {
			log.Println(err)
		}
	case "all":
		text := fetchAll(update.CallbackQuery.Data)
		err := bot.EditMessage(update.CallbackQuery.Message, text)
		if err != nil {
			log.Println(err)
		}
	case "docID":
		text := fetchOne(update.CallbackQuery.Data)
		err := bot.SendMessage(text, update.CallbackQuery.Message.Chat)
		if err != nil {
			log.Println(err)
		}

	case "search":
		bot.DeleteKeyboard()
		text := "Type Your Query: "
		newReply := query{UserID: update.CallbackQuery.From.ID,
			ChatID:    update.CallbackQuery.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.MessageID,
			Type:      "search"}
		replies = add(replies, newReply)
		err := bot.EditMessage(update.CallbackQuery.Message, text)
		if err != nil {
			log.Println(err)
		}

	case "contsearch":
		bot.DeleteKeyboard()
		newResult := query{
			UserID:    update.CallbackQuery.From.ID,
			MessageID: update.CallbackQuery.Message.MessageID,
			ChatID:    update.CallbackQuery.Message.Chat.ID,
			Text:      update.CallbackQuery.Data,
		}
		search(newResult)

	case "login":
		bot.DeleteKeyboard()
		text := "Please Enter Your Bookateria Account Credentials in the format below\n\nEmail:\nPassword:\n\n" +
			"Example: \njohndoe@gmail.com\nadmin@123!"
		newReply := query{UserID: update.CallbackQuery.From.ID,
			ChatID:    update.CallbackQuery.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.MessageID,
			Type:      "login"}

		replies = add(replies, newReply)
		err := bot.EditMessage(update.CallbackQuery.Message, text)
		if err != nil {
			log.Println(err)
		}
	}
}
