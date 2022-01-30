package main

const (
	urlContentSearch string = "https://mandrillapp.com/api/1.0/messages/search.json"
	urlContentInfo   string = "https://mandrillapp.com/api/1.0/messages/content.json"
	urlContentRemove string = "https://mandrillapp.com/api/1.0/rejects/delete"
)

type Payload struct {
	Key      string `json:"key"`
	ID       string `json:"id"`
	Query    string `json:"query"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
	Limit    int    `json:"limit"`
}

// var payload Payload

type returnContentSearch struct {
	ID      string `json:"_id"`
	Subject string `json:"subject"`
	Email   string `json:"email"`
	State   string `json:"state"`
	Opens   int    `json:"opens"`
	Clicks  int    `json:"clicks"`
	Ts      int64  `json:"ts"`
}

type returnContentInfo struct {
	ID      string `json:"_id"`
	Subject string `json:"subject"`
	To      struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"to"`
	TS   int64  `json:"ts"`
	Text string `json:"text"`
}
