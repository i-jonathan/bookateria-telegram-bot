package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/yoruba-codigy/goTelegram"
)

func fetchAll(callbackCode string) string {
	page := strings.Split(callbackCode, "-")[1]
	bot.DeleteKeyboard()
	var response ResponseStruct

	url := apiURL + "document?page=" + page + "&page_size=10"

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

This function handles document searching based on the text sent by the user
*/
func search(query query) {

	//remove query from list of awaiting replies
	defer func() {
		replies = remove(replies, &query)
	}()

	var (
		update   goTelegram.Update
		response ResponseStruct
	)

	params := url.Values{}
	page := 1
	pageSize := 8
	searchText := query.Text

	//check if this function is called from a prev or next button
	if strings.HasPrefix(query.Text, "contsearch") {
		parts := strings.Split(query.Text, "-")
		searchText = parts[1]
		page, _ = strconv.Atoi(parts[2])
	}

	//Initialize the chat to send document search results to
	update.Message.MessageID = query.MessageID
	update.Message.Chat.ID = query.ChatID

	params.Add("page", strconv.Itoa(page))
	params.Add("page_size", strconv.Itoa(pageSize))
	params.Add("search", searchText)

	//Make request to the api for documents with the specified title
	resp, err := http.Get(apiURL + "document?" + params.Encode())

	if err != nil {
		log.Println(err)
		go giveFeedback(&query, true)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		log.Println(err)
		go giveFeedback(&query, true)
		return
	}

	//If no document with the specified title is returned by the api, inform the user
	if len(response.Result) == 0 {
		bot.AddButton("Back", "documents")
		bot.MakeKeyboard(1)
		if err := bot.EditMessage(update.Message, "Couldn't Find Any Document That Matches Your Search Term: "+searchText); err != nil {
			log.Println(err)
		}
		return
	}

	//Process returned results from the api
	text := fmt.Sprintf("Showing Results For: %s\n", searchText)
	currIndex := (pageSize * page) - pageSize
	for _, doc := range response.Result {
		currIndex++
		text += fmt.Sprintf("%d. %s by %s\n", currIndex, doc.Title, doc.Author)
		bot.AddButton(strconv.Itoa(currIndex), "docID-"+strconv.Itoa(doc.Id))
	}

	bot.MakeKeyboard(len(response.Result))

	//Add prev and next buttons based on results returned from the api
	col := 0

	if response.Previous == true {
		bot.AddButton("Prev", "contsearch-"+searchText+"-"+strconv.Itoa(page-1))
		col += 1
	}

	if response.Next == true {
		bot.AddButton("Next", "contsearch-"+searchText+"-"+strconv.Itoa(page+1))
		col += 1
	}

	if col != 0 {
		bot.MakeKeyboard(col)
	}

	//Add button to go back to the documents Menu
	//if e.g user doesn't find desired document and wants to try another search keyword
	bot.AddButton("Documents Menu", "documents")
	bot.MakeKeyboard(1)

	if err := bot.EditMessage(update.Message, text); err != nil {
		log.Println(err)
	}
}

func giveFeedback(reply *query, err bool) {
	var text string
	var msg goTelegram.Update

	msg.Message.MessageID = reply.MessageID
	msg.Message.Chat.ID = reply.ChatID

	if err {
		text = "There Was An Error Processing Your Request, Please Try Again Later"
	} else {
		text = "Operation Completed Successfully"
	}

	bot.AddButton("Back", "documents")
	bot.MakeKeyboard(1)
	bot.EditMessage(msg.Message, text)
}

func processDocument(reply *query) {
	defer bot.DeleteKeyboard()
	var msg goTelegram.Update
	doc := mockDocs[reply.UserID]

	msg.Message.MessageID = reply.MessageID
	msg.Message.Chat.ID = reply.ChatID

	fileDets := reply.Update.Message.File

	text := "Upload The File Associated With The Details Provided In The Last Step"

	if fileDets.FileName != "" {
		doc.Filename = fileDets.FileName
		text += "\n\nFIle Received: " + fileDets.FileName
		bot.AddButton("Upload", "upload")
	}

	bot.AddButton("Cancel", "bail")
	bot.MakeKeyboard(1)
	log.Println(bot.EditMessage(msg.Message, text))
}

func fillDocument(reply *query) {
	defer bot.DeleteKeyboard()
	var (
		messageUpdate goTelegram.Update
		text          string
	)
	messageUpdate.Message.MessageID = reply.MessageID
	messageUpdate.Message.Chat.ID = reply.ChatID

	doc := mockDocs[reply.UserID]

	if reply.isEdit {
		switch reply.ReplyID {
		case doc.Title.messageID:
			doc.Title.text = reply.Text
		case doc.Author.messageID:
			doc.Author.text = reply.Text
		case doc.Summary.messageID:
			doc.Summary.text = reply.Text
		}
	} else {

		switch reply.Level {
		case 0:
			reply.Level++
		case 1:
			doc.Title.text = reply.Text
			doc.Title.messageID = reply.ReplyID
			reply.Level++

		case 2:
			doc.Author.text = reply.Text
			doc.Author.messageID = reply.ReplyID
			reply.Level++

		case 3:
			doc.Summary.text = reply.Text
			doc.Summary.messageID = reply.ReplyID
			reply.Level++

		}
	}

	text = fmt.Sprintf("Complete The Following Details:\n\nDocument Title: %s\nDocument Author: %s\nDocument Summary: %s\n", doc.Title.text, doc.Author.text, doc.Summary.text)
	if doc.Summary.text != "" {
		text = text + "\nPlease Check The Document Details And Verify They Are Correct, Press OK To Proceed"
		bot.AddButton("OK", "processDoc")
	}
	bot.AddButton("Cancel", "bail")
	bot.MakeKeyboard(1)

	if err := bot.EditMessage(messageUpdate.Message, text); err != nil {
		log.Println(err)
	}
}

func uploadDocument(chatID int) {

	var msg goTelegram.Update
	mockDoc := mockDocs[chatID]
	reply, _ := get(replies, chatID)
	fileDets := reply.Update.Message.File

	defer func() {
		delete(mockDocs, chatID)
		replies = remove(replies, reply)
	}()

	userToken, _ := getToken(chatID)

	msg.Message.MessageID = reply.MessageID
	msg.Message.Chat.ID = reply.ChatID

	if err := bot.GetFile(fileDets.FileID, fileDets.FileName); err != nil {
		log.Println("Couldn't Pull File From Telegram's Servers")
		log.Println(err)
		go giveFeedback(reply, true)
		return
	}

	rawFile, err := os.Open(mockDoc.Filename)
	if err != nil {
		log.Println("Couldn't Open " + mockDoc.Filename + " For Reading")
		log.Println(err)
		return
	}

	defer rawFile.Close()

	fileBody := new(bytes.Buffer)

	writer := multipart.NewWriter(fileBody)

	filePart, err := writer.CreateFormFile("file", mockDoc.Filename)

	io.Copy(filePart, rawFile)

	writer.WriteField("title", mockDoc.Title.text)
	writer.WriteField("author", mockDoc.Author.text)
	writer.WriteField("summary", mockDoc.Summary.text)

	writer.Close()

	req, _ := http.NewRequest("POST", apiURL+"document", fileBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", userToken)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		go giveFeedback(reply, true)
		log.Println("Couldn't Establish Connection With Bookateria's Servers")
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		go giveFeedback(reply, true)
		log.Println("Couldn't Upload File To Bookateria's Servers")
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(body))
		return
	}

	go giveFeedback(reply, false)

}
