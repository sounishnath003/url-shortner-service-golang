package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

// GenerateUrlShortenerHandler (v1) for generating the shorten url from the body provided.
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

	encodedUrl, err := GetMd5Hash(&url)
	if err != nil {
		handlers.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Convert these bytes to decimal: 1b3aabf5266b (hexadecimal) → 47770830013755 (decimal).
	num, err := strconv.ParseInt(encodedUrl, 16, 64)
	if err != nil {
		handlers.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	shortUrl, err := EncodeToBase62(num)
	if err != nil {
		handlers.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Grab core from context.
	co := r.Context().Value("co").(*core.Core)
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

// GetShortenUrlHandler (v1) gets the shorten url from the url provided in the path param.
func GetShortenUrlHandler(w http.ResponseWriter, r *http.Request) {
	// Grab the shortUrl from param
	shortUrl := r.URL.Path[len("/api/v1/"):]

	// Grab core from context.
	co := r.Context().Value("co").(*core.Core)

	_, err := co.QueryStmts.IncrUrlHitCountQuery.Exec(shortUrl)
	if err != nil {
		handlers.WriteError(w, http.StatusNotFound, errors.New("No url found for the given shorten url"))
		return
	}

	// Get the original url from the database.
	var originalUrl string
	var expirationAt time.Time

	err = co.QueryStmts.GetShortUrlQuery.QueryRow(shortUrl).Scan(&originalUrl, &expirationAt)
	if err != nil {
		handlers.WriteError(w, http.StatusNotFound, errors.New("No url found for the given shorten url"))
		return
	}

	if len(originalUrl) == 0 {
		handlers.WriteError(w, http.StatusBadGateway, errors.New("short url has been expired. or url found."))
	}

	// Redirect to the original url.
	http.Redirect(w, r, originalUrl, http.StatusFound)
}
