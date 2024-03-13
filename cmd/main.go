package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	shortenedLinks map[string]string // Maps shortened paths to original URLs.
)

func init() {
	// Seed the random number generator to ensure different results across runs.
	rand.Seed(time.Now().UnixNano())
}

func main() {
	shortenedLinks = make(map[string]string)

	// Route setup
	http.HandleFunc("/", handleHomePage)
	http.HandleFunc("/addLink", handleAddLink)
	http.HandleFunc("/short/", handleRedirectToOriginalLink)

	log.Println("Server is running on port :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}

// handleHomePage serves the home page and lists all shortened links.
func handleHomePage(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "text/html")
	writer.WriteHeader(http.StatusOK)

	var linksListHtml string
	for short, original := range shortenedLinks {
		linksListHtml += fmt.Sprintf("<div>Shortened: <a href=\"/short/%s\">/short/%s</a> Original: %s</div>", short, short, original)
	}

	fmt.Fprintf(writer, "<h2>Go URL Shortener</h2><p>Welcome! Here's a list of all shortened URLs:</p>%s", linksListHtml)
}

// handleAddLink creates a shortened link for the given original URL.
func handleAddLink(writer http.ResponseWriter, request *http.Request) {
	urlQuery := request.URL.Query()
	originalURLs, hasLink := urlQuery["link"]

	if !hasLink || !isValidUrl(originalURLs[0]) {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, "Invalid request. Please use an absolute URL, e.g., /addLink?link=https://example.com")
		return
	}

	originalURL := originalURLs[0]
	if _, exists := shortenedLinks[originalURL]; !exists {
		shortened := generateRandomString(10)
		shortenedLinks[shortened] = originalURL
		writer.Header().Set("Content-Type", "text/html")
		writer.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(writer, "Shortened URL: <a href=\"/short/%s\">/short/%s</a>", shortened, shortened)
	} else {
		writer.WriteHeader(http.StatusConflict)
		fmt.Fprint(writer, "This URL has already been shortened.")
	}
}

// isValidUrl checks if the provided URL is valid and absolute.
func isValidUrl(url string) bool {
	regex, err := regexp.Compile("^(http|https)://")
	if err != nil {
		log.Println("Failed to compile URL validation regex:", err)
		return false
	}
	return regex.MatchString(strings.TrimSpace(url))
}

// handleRedirectToOriginalLink redirects to the original URL based on the shortened identifier.
func handleRedirectToOriginalLink(writer http.ResponseWriter, request *http.Request) {
	pathComponents := strings.Split(request.URL.Path, "/")
	shortened := pathComponents[2]

	if originalURL, exists := shortenedLinks[shortened]; exists {
		http.Redirect(writer, request, originalURL, http.StatusTemporaryRedirect)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

// generateRandomString produces a random string of the specified length.
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = charset[rand.Intn(len(charset))]
	}
	return string(bytes)
}
