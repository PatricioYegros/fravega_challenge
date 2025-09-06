package models

type Event struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Date string `json:"date"`
	User string `json:"user"`
}
