package v1

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/sounishnath003/url-shortner-service-golang/internal/core"
	"github.com/sounishnath003/url-shortner-service-golang/internal/handlers"
)

// HealthHandler works as a health check endpoint for the api.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	co := r.Context().Value("co").(*core.Core)

	hostname, err := os.Hostname()
	// Handle err.
	if err != nil {
		handlers.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	handlers.WriteJson(w, http.StatusOK, map[string]any{
		"status":    "OK",
		"version":   co.Version,
		"message":   "api services are normal",
		"hostname":  hostname,
		"timestamp": time.Now(),
	})
}

// GenerateUrlShortenerV1Handler (v1) for generating the shorten url from the body provided.
func GenerateUrlShortenerV1Handler(w http.ResponseWriter, r *http.Request) {
	// Grab from body.
	var url CreateUShortenUrlDto
	json.NewDecoder(r.Body).Decode(&url)
	defer r.Body.Close()

	// Apply sanitize checks on the url.
	err := SanitizeURLChecks(&url)
	if err != nil {
		handlers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	short, err := GetMd5Hash(&url)
	if err != nil {
		handlers.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	handlers.WriteJson(w, http.StatusOK, map[string]any{
		"status":   "OK",
		"shortUrl": short,
		"O":        url,
		"message":  "short url generated",
	})
}

// GetShortenUrlV1Handler (v1) gets the shorten url from the url provided in the path param.
func GetShortenUrlV1Handler(w http.ResponseWriter, r *http.Request) {
	handlers.WriteJson(w, http.StatusOK, "Get shorten url")
}
