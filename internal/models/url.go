package models

import "time"

type Url struct {
	ID           int       `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	Hits         int       `json:"hits"`
	UserID       int       `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	ExpirationAt time.Time `json:"expiration_at"`
}
