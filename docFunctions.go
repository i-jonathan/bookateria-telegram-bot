package main

import (
	"encoding/json"
	"fmt"
	"github.com/yoruba-codigy/goTelegram"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func fetchAll(callbackCode string) string {
	page := strings.Split(callbackCode, "-")[1]
	bot.DeleteKeyboard()
	var response ResponseStruct

	url := apiURL + "document?page=" + page + "&page_size=1"

	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
		return ""
	}

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		log.Println(err)
		return ""
	}

	if len(response.Result) == 0 {
		return "No documents available"
	}

	text := "Documents available:\n"

	// Format text to be displayed in message

	for _, doc := range response.Result {
		text += fmt.Sprintf("%d. %s by %s\n", doc.Id, doc.Title, doc.Author)
		bot.AddButton(strconv.Itoa(doc.Id), "docID-"+strconv.Itoa(doc.Id))
	}

	bot.MakeKeyboard(len(response.Result))

	// Add button to go back
	bot.AddButton("Back", "documents")

	// Check if there is need for a previous page
	col := 0
	if response.Previous {
		bot.AddButton("Prev", "all-"+strconv.Itoa(response.Page-1))
		col += 1
	}

	if response.Next {
		bot.AddButton("Next", "all-"+strconv.Itoa(response.Page+1))
		col += 1
	}

	if col != 0 {
		bot.MakeKeyboard(col)
	}

	return text
}

func fetchOne(callbackCode string) string {
	bot.DeleteKeyboard()
	id := strings.Split(callbackCode, "-")[1]

	resp, err := http.Get(apiURL + "document/" + id)

	if err != nil {
		log.Println(err)
		return "An error occurred"
	}

	var doc document

	err = json.NewDecoder(resp.Body).Decode(&doc)
	if err != nil {
		log.Println(err)
		return "An error occurred"
	}

	text := fmt.Sprintf("Title: %s\nAuthor: %s\nDescription: %s\nEdition: %d\n\n %s",
		doc.Title, doc.Author, doc.Summary, doc.Edition, doc.FileSlug)

	return text
}

/*func search
Params: query
returns: nil

This funcntion handles document searching based on the text sent by the user
*/
func search(query query) {

	//remove query from list of awaiting replies
	defer func() {
		replies = remove(replies, query)
	}()

	var update goTelegram.Update
	var response ResponseStruct

	page := 1
	search_text := query.Text

	//check if this function is called from a prev or next button
	if strings.HasPrefix(query.Text, "contsearch") {
		parts := strings.Split(query.Text, "-")
		search_text = parts[1]
		page, _ = strconv.Atoi(parts[2])
	}

	//Initialize the chat to send document search results to
	update.Message.MessageID = query.Message_ID
	update.Message.Chat.ID = query.Chat_ID

	//Make request to the api for documents with the specified title
	resp, err := http.Get(apiURL + "document?page_size=1&search=" + search_text + "&page=" + strconv.Itoa(page))

	if err != nil {
		log.Println(err)
		bot.AddButton("Back", "documents")
		bot.MakeKeyboard(1)
		bot.EditMessage(update.Message, "An Error Occured While Processing Your Request")
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		log.Println(err)
		bot.AddButton("Back", "documents")
		bot.MakeKeyboard(1)
		bot.EditMessage(update.Message, "Couldn't Find Any Document That Matches Your Search Term: "+search_text)
		return
	}

	//If no document with the specified title is returned by the api, inform the user
	if len(response.Result) == 0 {
		bot.AddButton("Back", "documents")
		bot.MakeKeyboard(1)
		bot.EditMessage(update.Message, "Couldn't Find Any Document That Matches Your Search Term")
		return
	}

	//Process returned results from the api
	text := fmt.Sprintf("Showing Results For: %s\n", search_text)

	for index, doc := range response.Result {
		text += fmt.Sprintf("%d. %s by %s\n", index+1, doc.Title, doc.Author)
		bot.AddButton(strconv.Itoa(index+1), "docID-"+strconv.Itoa(doc.Id))
	}

	bot.MakeKeyboard(len(response.Result))

	//Add prev and nexr buttons based on results returned from the api
	col := 0

	if response.Previous == true {
		bot.AddButton("Prev", "contsearch-"+search_text+"-"+strconv.Itoa(page-1))
		col += 1
	}

	if response.Next == true {
		bot.AddButton("Next", "contsearch-"+search_text+"-"+strconv.Itoa(page+1))
		col += 1
	}

	if col != 0 {
		bot.MakeKeyboard(col)
	}

	//Add button to go back to the documents Menu
	//if e.g user doesn't find desired document and wants to try another search keyword
	bot.AddButton("Douments Menu", "documents")
	bot.MakeKeyboard(1)

	bot.EditMessage(update.Message, text)
}
