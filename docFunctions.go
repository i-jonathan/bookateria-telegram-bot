package main

import (
	"encoding/json"
	"fmt"
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
