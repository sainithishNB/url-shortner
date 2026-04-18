package models

import "time"

type URL struct {
	ID         uint
	ShortCode  string `gorm:"size:10;uniqueIndex"`
	LongURL    string `gorm:"type:text"`
	CreatedAt  time.Time
	ExpiresAt  *time.Time
	ClickCount int
}
type ShortenRequest struct {
	URL       string `json:"url"`
	Alias     string `json:"alias"`
	ExpiresIn int    `json:"expires_in"`
}
