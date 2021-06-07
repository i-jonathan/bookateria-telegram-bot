package main

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

type query struct {
	UserID    int    //UserID
	ChatID    int    //ChatID Of The Incoming Text
	MessageID int    //Message ID Of The Text To Be Updated
	Level     int    //Question Level For Multi-Level Questions
	Type      string //Query Type e.g Search, Login e.t.c
	Text      string //Text Of The Incoming update
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
