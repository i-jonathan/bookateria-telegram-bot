package main

type document struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Edition  int    `json:"edition"`
	Author   string `json:"author"`
	Summary  string `json:"summary"`
	FileSlug string `json:"file_slug"`
}