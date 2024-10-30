package v2

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"
)

var (
	BASE_CHARACTERS = "abcdefghijklmnopqrsuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	BASE_LEN        = len(BASE_CHARACTERS)
)

// SanitizeURLChecks helps to sanitize the url before the creation
// shorten urls.
// This will also fill the default expiry to parameter if the expiry date is not provided.
func SanitizeURLChecks(urlInfo *CreateUShortenUrlDto) error {
	// Length check.
	if len(urlInfo.OriginalUrl) < 5 {
		return fmt.Errorf("url is too short: %s", urlInfo.OriginalUrl)
	}

	// Check if a valid url.
	urlScheme, err := url.ParseRequestURI(urlInfo.OriginalUrl)
	if err != nil {
		return err
	}

	// Check the URL scheme to be only Http or https
	if urlScheme.Scheme == "" || !(urlScheme.Scheme == "http" ||
		urlScheme.Scheme == "https") {
		return fmt.Errorf("invalid url scheme: %s", urlInfo.OriginalUrl)
	}

	// Add default expiration - 2 DAY default.
	if urlInfo.ExpiryDate.Before(time.Now()) {
		urlInfo.ExpiryDate = time.Now().Add(48 * time.Hour)
	}

	return nil
}

// GetMd5Hash helps to generate the 6 bytes hash encoded infromation.
func GetMd5Hash(urlInfo *CreateUShortenUrlDto) (string, error) {
	hasher := md5.New()
	hasher.Write([]byte(urlInfo.OriginalUrl))
	shortenUrl := hex.EncodeToString(hasher.Sum(nil))[:6]
	return shortenUrl, nil
}

// EncodeToBase62 helps to generate the base62 encoded string from the number int64.
//
// Encode the result into a Base62 encoded string: DZFbb43.
func EncodeToBase62(num int64) (string, error) {
	shortUrl := make([]byte, 1)
	k := 0
	for num > 0 {
		rem := num % int64(BASE_LEN)
		shortUrl = append(shortUrl, []byte(BASE_CHARACTERS)[rem])
		num = num / 10
		k += 1
	}

	return string(shortUrl[1:]), nil
}
