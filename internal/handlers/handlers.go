package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/sounishnath003/url-shortner-service-golang/internal/core"
)

// GetShortenUrlHandler gets the shorten url from the url provided in the path param.
func GetShortenUrlHandler(w http.ResponseWriter, r *http.Request) {
	// Grab the shortUrl from param
	shortUrl := r.URL.Path[1:]

	// Grab core from context.
	co := r.Context().Value("co").(*core.Core)

	_, err := co.QueryStmts.IncrUrlHitCountQuery.Exec(shortUrl)
	if err != nil {
		WriteError(w, http.StatusNotFound, errors.New("No url found for the given shorten url"))
		return
	}

	// Check the URL is present in cache.
	originalUrl, err := co.FindOriginalUrlFromCache(shortUrl)
	if err == nil && len(originalUrl) > 0 {
		co.Lo.Info("[CACHE_HIT]", "originalUrl", originalUrl, "shortUrl", shortUrl)
		http.Redirect(w, r, originalUrl, http.StatusFound)
		return
	}
	co.Lo.Info("[CACHE_MISS]", "originalUrl", originalUrl, "shortUrl", shortUrl)

	// Get the original url from the database.
	var expirationAt time.Time

	err = co.QueryStmts.GetShortUrlQuery.QueryRow(shortUrl).Scan(&originalUrl, &expirationAt)
	if err != nil {
		WriteError(w, http.StatusNotFound, errors.New("No url found for the given shorten url"))
		return
	}

	if len(originalUrl) == 0 {
		WriteError(w, http.StatusBadGateway, errors.New("short url has been expired. or url found."))
	}

	// Redirect to the original url.
	http.Redirect(w, r, originalUrl, http.StatusFound)
}

func CustomAliasAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	// Grab the alias.
	customAlias := r.URL.Path[len("/api/check-alias/"):]

	// Get the core.Core context
	co := r.Context().Value("co").(*core.Core)
	// Check the existence.
	co.Lo.Info("checking customAlias available using bloom filter", "alias", customAlias)
	_, exists := co.BloomFilter.Exists(customAlias)
	if exists {
		WriteJson(w, http.StatusNotAcceptable, "alias is already taken. Try another one.")
		return
	}

	WriteJson(w, http.StatusOK, "alias is available.")
}
