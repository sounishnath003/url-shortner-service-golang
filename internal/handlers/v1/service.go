package v1

import (
	"crypto/md5"
)

// GetMd5Hash helps to generate the encode infromation
func GetMd5Hash(urlInfo *CreateUShortenUrlDto) ([]byte, error) {
	url := urlInfo.LongUrl
	urlBytes := md5.New().Sum([]byte(url))[:6]
	return urlBytes, nil
}
