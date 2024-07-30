package models

type Stats struct {
	TotalMessages int `json:"total_messages"`
	ProcessedMessages int `json:"processed_messages"`
}

type Message struct {
	ID int `json:"id"`
	Content string `json:"content"`
	Processed bool `json:"processed"`
}