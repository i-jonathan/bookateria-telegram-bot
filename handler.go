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
		var buttons []goTelegram.InlineKeyboard
		bot.DeleteKeyboard()

		if update.Message.Chat.Type == "private" {
			buttons = append(buttons, goTelegram.InlineKeyboard{
				Text: "Accounts",
				Data: "accounts",
			})
		}

		buttons = append(buttons, goTelegram.InlineKeyboard{
			Text: "Documents",
			Data: "documents",
		})

		bot.MakeKeyboard(buttons, 2)
		bot.SendMessage("Hello. Where would you like to go?", chat)
	}
}

func processCallback(update goTelegram.Update) {
	//chat := update.Message.Chat
	defer bot.AnswerCallback(update.CallbackQuery.ID)

	switch update.CallbackQuery.Data {
	case "documents":
		var buttons []goTelegram.InlineKeyboard
		bot.DeleteKeyboard()

		if update.Message.Chat.Type == "private" {
			buttons = append(buttons, goTelegram.InlineKeyboard{
				Text: "Add",
				Data: "add",
			})

			buttons = append(buttons, goTelegram.InlineKeyboard{
				Text: "Update",
				Data: "update",
			})

			buttons = append(buttons, goTelegram.InlineKeyboard{
				Text: "My Uploads",
				Data: "mine",
			})
		}

		buttons = append(buttons, goTelegram.InlineKeyboard{
			Text: "All",
			Data: "all",
		})

		buttons = append(buttons, goTelegram.InlineKeyboard{
			Text: "Search",
			Data: "search",
		})

		buttons = append(buttons, goTelegram.InlineKeyboard{
			Text: "Tags",
			Data: "tags",
		})

		buttons = append(buttons, goTelegram.InlineKeyboard{
			Text: "Categories",
			Data: "cat",
		})

		bot.MakeKeyboard(buttons, 3)
		bot.EditMessage(update.CallbackQuery.Message, "Documents")
	}
}