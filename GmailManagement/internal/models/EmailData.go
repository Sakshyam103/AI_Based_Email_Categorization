package models

type Attachment struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
	MimeType string `json:"mimetype"`
}

type EmailData struct {
	To          string       `json:"to"`
	Subject     string       `json:"subject"`
	Body        string       `json:"body"`
	Attachments []Attachment `json:"attachments"`
}
