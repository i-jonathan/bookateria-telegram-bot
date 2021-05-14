package main

type ResponseStruct struct {
	Previous bool        `json:"previous"`
	Next     bool        `json:"next"`
	Page     int         `json:"page"`
	Count    int64       `json:"count"`
	Result   []document `json:"result"`
}

type document struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Edition  int    `json:"edition"`
	Author   string `json:"author"`
	Summary  string `json:"summary"`
	FileSlug string `json:"file_slug"`
}

type query struct {
	User       int
	Chat_ID    int
	Message_ID int
	Text       string
}
