package models

type CategorizedEmail struct {
	EmailId   string `json:"email_id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	HistoryId uint64 `json:"history_id"`
	Body      string `json:"body"`
	Date      string `json:"date"`
	Category  string `json:"category"`
}
