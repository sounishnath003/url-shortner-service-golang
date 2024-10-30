package v2

import "time"

type CreateUShortenUrlDto struct {
	OriginalUrl string    `json:"original_url"`
	CustomAlias string    `json:"custom_alias"`
	ExpiryDate  time.Time `json:"expiry_date"`
}
