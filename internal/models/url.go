package models

import "time"

type Url struct {
	ID           int       `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	Hits         int       `json:"hits"`
	CreatedAt    time.Time `json:"created_at"`
	ExpirationAr time.Time `json:"expiration_at"`
}
