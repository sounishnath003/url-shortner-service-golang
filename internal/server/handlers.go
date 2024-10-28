package server

import (
	"net/http"
	"os"
	"time"
)

// HealthHandler works as a health check endpoint for the api.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
	}

	WriteJson(w, http.StatusOK, map[string]any{
		"status":    "OK",
		"message":   "api services are normal",
		"hostname":  hostname,
		"timestamp": time.Now(),
	})
}

// GenerateUrlShortenerV1Handler for generating the shorten url from the body provided.
func GenerateUrlShortenerV1Handler(w http.ResponseWriter, r *http.Request) {
	WriteJson(w, http.StatusOK, "URL has been shorten")
}

// GetShortenUrlV1Handler gets the shorten url from the url provided in the path param.
func GetShortenUrlV1Handler(w http.ResponseWriter, r *http.Request) {
	WriteJson(w, http.StatusOK, "Get shorten url")
}
