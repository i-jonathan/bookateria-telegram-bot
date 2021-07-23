package main

import "github.com/yoruba-codigy/goTelegram"

type ResponseStruct struct {
	Previous bool       `json:"previous"`
	Next     bool       `json:"next"`
	Page     int        `json:"page"`
	Count    int64      `json:"count"`
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

type field struct {
	text      string
	messageID int
}

type mockDocument struct {
	Title    field
	Edition  field
	Author   field
	Summary  field
	Filename string
}

type query struct {
	UserID    int               //ID Of The Entity That Sent The Text
	ChatID    int               //ChatID Of The Incoming Text
	ReplyID   int               //ID Of The Incoming Message
	MessageID int               //Message ID Of The Text To Be Updated
	Level     int               //Question Level For Multi-Level Questions
	Type      string            //Query Type e.g Search, Login e.t.c
	Text      string            //Text Of The Incoming update
	Update    goTelegram.Update //Might Wanna Access Other Things From The Incoming Messae Body
	isEdit    bool              //If Incoming Reply Is An Edit Of A Previous Message
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	StayIn   bool   `json:"stay_in"`
}

type LoginResponse struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Expiry string `json:"expiry"`
}
