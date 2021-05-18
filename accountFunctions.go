package main

import "github.com/yoruba-codigy/goTelegram"

func handleLogin(update query) {
	//Do Random Login Stuff Here

	defer func() {
		replies = remove(replies, update)
	}()

	var replyUpdate goTelegram.Update

	replyUpdate.Message.MessageID = update.Message_ID
	replyUpdate.Message.Chat.ID = update.Chat_ID

	bot.EditMessage(replyUpdate.Message, "Login Code Should Be Here")
}
