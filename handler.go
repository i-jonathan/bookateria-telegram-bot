package main

import (
	"log"
	"strings"

	"github.com/yoruba-codigy/goTelegram"
)

func handler(update goTelegram.Update) {
	checkReplies(update)

	switch update.Type {
	case "text":
		go processRequest(update)
	case "callback":
		go processCallback(update)
	case "edited_text":
		go processEditedText(update)
	}
}

func checkReplies(update goTelegram.Update) {

	if update.Type == "edited_text" {
		update.Message = update.EditedMessage
	}
	if reply, ok := get(replies, update.Message.Chat.ID); ok {
		reply.Update = update
		reply.isEdit = false
		reply.ReplyID = update.Message.MessageID
		if update.Type == "edited_text" {
			reply.isEdit = true
		}

		if reply.UserID == update.Message.From.ID {
			reply.Text = update.Message.Text
			switch reply.Type {
			case "search":
				search(*reply)
			case "login":
				handleLogin(*reply)
			case "upload":
				fillDocument(reply)
			case "processDoc":
				processDocument(reply)
			}
		}
	}

}

func showMainMenu(update goTelegram.Update) {
	bot.DeleteKeyboard()
	var msg goTelegram.Update

	text := "Hello, Where Would You Like To Go?"

	if update.Message.Chat.Type == "private" || update.CallbackQuery.Message.Chat.Type == "private" {
		bot.AddButton("Account", "account")
	}
	bot.AddButton("Documents", "documents")
	bot.MakeKeyboard(2)

	if update.Type == "text" {
		chat := update.Message.Chat
		if err := bot.SendMessage(text, chat); err != nil {
			log.Println(err)
		}
	} else {
		msg.Message.MessageID = update.CallbackQuery.Message.MessageID
		msg.Message.Chat.ID = update.CallbackQuery.From.ID
		if err := bot.EditMessage(msg.Message, text); err != nil {
			log.Println(err)
		}
	}
}

func processRequest(update goTelegram.Update) {

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

		showMainMenu(update)

	}
}

func processEditedText(update goTelegram.Update) {
	upd := goTelegram.Update{
		Message: update.EditedMessage,
	}

	processRequest(upd)
}

func processCallback(update goTelegram.Update) {
	answeredCallback := false
	defer func() {
		if !answeredCallback {

			err := bot.AnswerCallback(update.CallbackQuery.ID, "", false)

			if err != nil {
				log.Println(err)
			}
		}
	}()

	command := update.CallbackQuery.Data

	if update.CallbackQuery.Data == "bail" {
		replies = removeByChatID(replies, update.CallbackQuery.From.ID)
		command = "documents"
	}

	if strings.HasPrefix(update.CallbackQuery.Data, "docID") {
		command = "docID"
	} else if strings.HasPrefix(update.CallbackQuery.Data, "all") {
		command = "all"
	} else if strings.HasPrefix(update.CallbackQuery.Data, "contsearch") {
		command = "contsearch"
	}

	switch command {

	//Menu Cases
	case "account":
		bot.DeleteKeyboard()
		bot.AddButton("Login", "login")
		bot.AddButton("Register", "register")
		bot.AddButton("Logout", "logout")
		bot.AddButton("Back", "main-menu")
		bot.MakeKeyboard(3)
		if err := bot.EditMessage(update.CallbackQuery.Message, "Accounts"); err != nil {
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
		bot.AddButton("Back", "main-menu")

		bot.MakeKeyboard(3)
		if err := bot.EditMessage(update.CallbackQuery.Message, "Documents"); err != nil {
			log.Println(err)
		}

	case "main-menu":
		showMainMenu(update)

		//Operation Cases
	case "all":
		text := fetchAll(update.CallbackQuery.Data)
		if err := bot.EditMessage(update.CallbackQuery.Message, text); err != nil {
			log.Println(err)
		}

	case "add":
		bot.DeleteKeyboard()

		if _, err := getToken(update.CallbackQuery.From.ID); err != nil {
			bot.AnswerCallback(update.CallbackQuery.ID, "Login Is Required To Add A Document", true)
			answeredCallback = true
			return
		}

		newReply := &query{
			UserID:    update.CallbackQuery.From.ID,
			ChatID:    update.CallbackQuery.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.MessageID,
			Type:      "upload",
		}

		mockDocs[update.CallbackQuery.From.ID] = &mockDocument{}
		replies = add(replies, newReply)
		fillDocument(newReply)

	case "docID":
		text := fetchOne(update.CallbackQuery.Data)
		if err := bot.SendMessage(text, update.CallbackQuery.Message.Chat); err != nil {
			log.Println(err)
		}

	case "search":
		bot.DeleteKeyboard()
		text := "Type Your Query: "
		newReply := &query{UserID: update.CallbackQuery.From.ID,
			ChatID:    update.CallbackQuery.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.MessageID,
			Type:      "search"}
		replies = add(replies, newReply)
		if err := bot.EditMessage(update.CallbackQuery.Message, text); err != nil {
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
		newReply := &query{UserID: update.CallbackQuery.From.ID,
			ChatID:    update.CallbackQuery.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.MessageID,
			Type:      "login"}

		replies = add(replies, newReply)
		bot.MakeKeyboard(1)
		if err := bot.EditMessage(update.CallbackQuery.Message, text); err != nil {
			log.Println(err)
		}

	case "processDoc":
		bot.DeleteKeyboard()
		reply, _ := get(replies, update.CallbackQuery.From.ID)
		reply.Type = "processDoc"
		processDocument(reply)

	case "upload":
		uploadDocument(update.CallbackQuery.From.ID)
		delete(mockDocs, update.CallbackQuery.From.ID)
	}
}
