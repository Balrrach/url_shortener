package main

import "time"

// ShortURL represents a shortened URL in our system
type ShortURL struct {
	ID        int64     `json:"id"`
	ShortCode string    `json:"short_code"`
	LongURL   string    `json:"long_url"`
	CreatedAt time.Time `json:"created_at"`
	Clicks    int64     `json:"clicks"`
}

// URLRequest is used to accept incoming URL shortening requests
type URLRequest struct {
	LongURL string `json:"url" binding:"required"`
}
