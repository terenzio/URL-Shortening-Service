package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/terenzio/URL-Shortening-Service/application"
)

type Handler struct {
	service *application.URLService
}

func NewHandler(service *application.URLService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleHomePage(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}
	writer.Header().Set("Content-Type", "text/html")
	urls, err := h.service.FetchAllURLs(request.Context())
	if err != nil {
		http.Error(writer, "Failed to fetch URLs", http.StatusInternalServerError)
		return
	}

	var linksListHtml string
	for _, url := range urls {
		linksListHtml += fmt.Sprintf("<div>Shortened: <a href=\"/short/%s\">/short/%s</a> Original: %s</div>", url.ShortCode, url.ShortCode, url.OriginalURL)
	}

	fmt.Fprintf(writer, "<h2>Go URL Shortener</h2><p>Welcome! Here's a list of all shortened URLs:</p>%s", linksListHtml)
}

func (h *Handler) HandleAddLink(writer http.ResponseWriter, request *http.Request) {
	urlQuery := request.URL.Query()
	originalURLs, hasLink := urlQuery["link"]

	if !hasLink || !isValidUrl(originalURLs[0]) {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, "Invalid request. Please use an absolute URL, e.g., /addLink?link=https://example.com")
		return
	}
	originalURL := originalURLs[0]

	// Here we hardcode the expiry duration to 30 days.
	// Modify as needed or extract from the request for dynamic durations.
	expiryDuration := 30 * 24 * time.Hour
	shortCode, err := h.service.ShortenURL(request.Context(), originalURL, expiryDuration)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error shortening URL: %v", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "text/html")
	writer.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(writer, "<h2>Go URL Shortener</h2>")
	fmt.Fprintf(writer, "<p>Original URLs: %s<br>", originalURL)
	fmt.Fprintf(writer, "<p>Shortened URL: <a href=\"/short/%s\">/short/%s</a> </p>", shortCode, shortCode)
}

func (h *Handler) HandleRedirectToOriginalLink(writer http.ResponseWriter, request *http.Request) {
	shortCode := request.URL.Path[len("/short/"):]
	if shortCode == "" {
		http.NotFound(writer, request)
		return
	}

	originalURL, err := h.service.RedirectToOriginalURL(request.Context(), shortCode)
	if err != nil {
		http.Error(writer, "Shortened URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(writer, request, originalURL, http.StatusTemporaryRedirect)
}
