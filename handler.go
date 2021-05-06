package main

import (
	"github.com/yoruba-codigy/goTelegram"
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
		bot.SendMessage("Hello. Where would you like to go?", chat)
	}
}

func processCallback(update goTelegram.Update) {
	defer bot.AnswerCallback(update.CallbackQuery.ID)

	var command string

	if strings.HasPrefix(update.CallbackQuery.Data, "docID") {
		command = "docID"
	} else {
		command = update.CallbackQuery.Data
	}

	switch command {
	case "documents":
		bot.DeleteKeyboard()

		if update.CallbackQuery.Message.Chat.Type == "private" {
			bot.AddButton("Add", "add")
			bot.AddButton("Update", "update")
			bot.AddButton("My Uploads", "mine")
		}

		bot.AddButton("All", "all")
		bot.AddButton("Search", "search")
		bot.AddButton("Tags", "tags")
		bot.AddButton("Categories", "cat")


		bot.MakeKeyboard(3)
		bot.EditMessage(update.CallbackQuery.Message, "Documents")
	case "all":
		text := fetchAll("1")
		bot.EditMessage(update.CallbackQuery.Message, text)
	case "docID":
		text := fetchOne(update.CallbackQuery.Data)
		bot.SendMessage(text, update.CallbackQuery.Message.Chat)
	}
}
