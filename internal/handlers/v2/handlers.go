package v2

import (
	"encoding/json"
	"errors"
	"fmt"
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

// GenerateUrlShortenerHandler (v2) for generating the shorten url from the body provided.
// Uses the database as atomic and safe atomic incremental id generation.
// Better than the hashing based methods in /api/v1/*
func GenerateUrlShortenerHandler(w http.ResponseWriter, r *http.Request) {
	// Grab from body.
	var url CreateUShortenUrlDto
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		handlers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	defer r.Body.Close()

	// Apply sanitize checks on the url.
	err = SanitizeURLChecks(&url)
	if err != nil {
		handlers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Grab the from context.
	co := r.Context().Value("co").(*core.Core)

	// Check if custom alias provided in body.
	shortUrl := url.CustomAlias

	// If no shortUrl is defined then generate a unique short url
	if len(shortUrl) == 0 {
		var num int
		// Get the incremental ID (distributed-ACID-compliant) safe
		// Comes at cost of performance in read heavy environment.
		co.QueryStmts.GetIncrementalIDQuery.QueryRow().Scan(&num)

		shortUrl, err = EncodeToBase62(int64(num))
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		// Double check the alias
		_, exists := co.BloomFilter.Exists(shortUrl)
		if exists {
			handlers.WriteError(w, http.StatusNotAcceptable, errors.New("alias is already taken. Try another one."))
			return
		}
	}

	// Get the user from context
	userID := r.Context().Value("userID").(int)

	// Save it to database.
	err = co.CreateNewShortUrlAsTxn(url.OriginalUrl, shortUrl, url.ExpiryDate, userID)
	if err != nil {
		handlers.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	handlers.WriteJson(w, http.StatusOK, map[string]any{
		"shortUrl": fmt.Sprintf("%s/%s", r.Host, shortUrl),
		"message":  "short url has been generated",
		"expiryBy": url.ExpiryDate,
	})
}
