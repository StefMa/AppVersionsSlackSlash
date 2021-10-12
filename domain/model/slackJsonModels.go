package model

type SlackJsonRequest struct {
	Blocks []SlackBlock `json:"blocks"`
}

type SlackBlock struct {
	Type string    `json:"type"`
	Text SlackText `json:"text"`
}

type SlackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
