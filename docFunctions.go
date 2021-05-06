package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func fetchAll(page string) string {
	bot.DeleteKeyboard()
	var allDocs []document

	url := apiURL + "document?page=" + page + "page_size=5"

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

	for _, doc := range allDocs {
		text += fmt.Sprintf("%d. %s by %s\n", doc.Id, doc.Title, doc.Author)
		bot.AddButton(strconv.Itoa(doc.Id), "docID-" + strconv.Itoa(doc.Id))
	}

	bot.MakeKeyboard(len(allDocs))
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