package v1

import "time"

type CreateUShortenUrlDto struct {
	LongUrl     string    `json:"longUrl"`
	CustomAlias string    `json:"customAlias"`
	ExpiryDate  time.Time `json:"expiryDate"`
}