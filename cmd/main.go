package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const base62Characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

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
		shortened := generateShortCode(originalURL)
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
// If the shortened identifier is not found, the function returns a 404 Not Found status.
func handleRedirectToOriginalLink(writer http.ResponseWriter, request *http.Request) {
	pathComponents := strings.Split(request.URL.Path, "/")
	shortened := pathComponents[2]

	if originalURL, exists := shortenedLinks[shortened]; exists {
		http.Redirect(writer, request, originalURL, http.StatusTemporaryRedirect)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

// hashURL creates a SHA-256 hash of the given URL.
// It uses the sequence number to handle hash collisions and ensure uniqueness.
// The function returns a byte slice of the hash.
func hashURL(url string, sequence int) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s%d", url, sequence)))
	return hasher.Sum(nil)
}

// base62Encode encodes a byte slice into a Base62 string.
// It uses a big.Int to handle large numbers and a lookup table for the encoding.
// The function returns a string of up to 8 characters, which is the maximum length for a 64-bit integer.
func base62Encode(bytes []byte) string {
	var result []byte
	number := new(big.Int).SetBytes(bytes)
	base := big.NewInt(62)
	zero := big.NewInt(0)
	mod := new(big.Int)

	for number.Cmp(zero) != 0 {
		number.DivMod(number, base, mod)
		result = append(result, base62Characters[mod.Int64()])
	}

	// Ensure the result is 7 characters by padding with "0" (or another character),
	// if necessary. This is a simple form of error handling and may not suit all cases.
	for len(result) < 8 {
		result = append(result, '0')
	}

	// Reverse the result since the encoding process generates it in reverse order.
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

// isUnique checks if the given short code is unique.
// It returns true if the short code is unique, and false otherwise.
func isUnique(shortCode string) bool {
	_, exists := shortenedLinks[shortCode]
	return !exists
}

// generateShortCode creates a unique short code for the given URL.
// It uses a sequence number to handle hash collisions and ensure uniqueness.
// The function generates a SHA-256 hash of the URL and encodes it using Base62.
// The result is a short code of up to 8 characters, which is then checked for uniqueness.
// If the short code is not unique, the sequence number is incremented and the process is repeated.
func generateShortCode(url string) string {
	sequence := 1
	for {
		hash := hashURL(url, sequence)
		partialHash := hash[:10]
		shortCode := base62Encode(partialHash)

		if len(shortCode) > 8 {
			shortCode = shortCode[:8]
		}

		if isUnique(shortCode) {
			shortenedLinks[shortCode] = url // Store the unique short code and its original URL
			return shortCode
		}
		sequence++
	}
}
