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
	var allDocs []document

	url := apiURL + "document?page=" + page + "&page_size=1"

	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
		return ""
	}

	err = json.NewDecoder(resp.Body).Decode(&allDocs)

	if err != nil {
		log.Println(err)
		return ""
	}

	if len(allDocs) == 0 {
		return "No documents available"
	}

	text := "Documents available:\n"

	// Format text to be displayed in message

	for _, doc := range allDocs {
		text += fmt.Sprintf("%d. %s by %s\n", doc.Id, doc.Title, doc.Author)
		bot.AddButton(strconv.Itoa(doc.Id), "docID-" + strconv.Itoa(doc.Id))
	}

	bot.MakeKeyboard(len(allDocs))

	// Add button to go back
	bot.AddButton("Back", "documents")

	// Check if there is need for a previous page
	current, err := strconv.Atoi(page)
	if current - 1 != 0 {
		bot.AddButton("Prev", "all-" + strconv.Itoa(current-1))
	}

	// Check if there is need for a next page
	url = apiURL + "document?page=" + strconv.Itoa(current+1) + "&page_size=1"

	resp, err = http.Get(url)

	if err != nil {
		log.Println(err)
		return ""
	}

	err = json.NewDecoder(resp.Body).Decode(&allDocs)

	if err != nil {
		log.Println(err)
		return ""
	}

	if len(allDocs) != 0 {
		bot.AddButton("Next", "all-" + strconv.Itoa(current+1))
	}
	bot.MakeKeyboard(2)

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