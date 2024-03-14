package http

import (
	"log"
	"regexp"
	"strings"
)

// isValidUrl checks if the given URL is valid.
// It returns true if the URL is valid, and false otherwise.
// A valid URL is an absolute URL with either the http or https scheme.
// For example, https://example.com is a valid URL, while example.com is not.
func isValidUrl(url string) bool {
	regex, err := regexp.Compile("^(http|https)://")
	if err != nil {
		log.Println("Failed to compile URL validation regex:", err)
		return false
	}
	return regex.MatchString(strings.TrimSpace(url))
}
