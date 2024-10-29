package v1

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"time"
)

// SanitizeURLChecks helps to sanitize the url before the creation
// shorten urls.
func SanitizeURLChecks(urlInfo *CreateUShortenUrlDto) error {
	// Length check.
	if len(urlInfo.OriginalUrl) < 5 {
		return fmt.Errorf("url is too short: %s", urlInfo.OriginalUrl)
	}

	// Check if a valid url.
	_, err := url.ParseRequestURI(urlInfo.OriginalUrl)
	if err != nil {
		return err
	}

	// Add default expiration - 2 DAY default.
	if urlInfo.ExpiryDate.Before(time.Now()) {
		urlInfo.ExpiryDate = time.Now().Add(48 * time.Hour)
	}

	return nil
}

// GetMd5Hash helps to generate the encode infromation
func GetMd5Hash(urlInfo *CreateUShortenUrlDto) ([]byte, error) {
	url := urlInfo.OriginalUrl
	urlBytes := md5.New().Sum([]byte(url))[:6]
	return urlBytes, nil
}
