package server

import (
	"net/http"
	"os"
	"time"

	"github.com/sounishnath003/url-shortner-service-golang/internal/core"
)

// HealthHandler works as a health check endpoint for the api.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	co := r.Context().Value("co").(*core.Core)

	hostname, err := os.Hostname()
	// Handle err.
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	WriteJson(w, http.StatusOK, map[string]any{
		"status":    "OK",
		"version":   co.Version,
		"message":   "api services are normal",
		"hostname":  hostname,
		"timestamp": time.Now(),
	})
}

// GenerateUrlShortenerV1Handler (v1) for generating the shorten url from the body provided.
func GenerateUrlShortenerV1Handler(w http.ResponseWriter, r *http.Request) {
	WriteJson(w, http.StatusOK, "URL has been shorten")
}

// GetShortenUrlV1Handler (v1) gets the shorten url from the url provided in the path param.
func GetShortenUrlV1Handler(w http.ResponseWriter, r *http.Request) {
	WriteJson(w, http.StatusOK, "Get shorten url")
}
