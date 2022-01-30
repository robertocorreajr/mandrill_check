package entity

import "os"

// Payload Struct (Model)
type payload struct {
	Key      string `json:"key"`
	ID       string `json:"id"`
	Query    string `json:"query"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
	Limit    int    `json:"limit"`
	Email    string `json:"email"`
}

func NewPayload() *payload {
	payload := payload{}
	return &payload
}

type returnContentSearch struct {
	ID      string `json:"_id"`
	Subject string `json:"subject"`
	Email   string `json:"email"`
	State   string `json:"state"`
	Opens   int    `json:"opens"`
	Clicks  int    `json:"clicks"`
	Ts      int64  `json:"ts"`
}

func (p payload) Search(days int) (returnContentSearch, error) {
	if days == 0 {
		days = -7
	}

	p.Key = os.Getenv("KEY")

	contentBody := returnContentSearch{
		ID:      "1",
		Subject: "teste",
	}
	return contentBody, nil
}

// func (p payload) Info() (returnContentSearch, error) {

// }
