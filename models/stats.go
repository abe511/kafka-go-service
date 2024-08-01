package models

type Stats struct {
	TotalMessages int `json:"total_messages"`
	ProcessedMessages int `json:"processed_messages"`
}
