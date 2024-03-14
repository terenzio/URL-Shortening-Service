package http

import (
	"log"
	"regexp"
	"strings"
)

func isValidUrl(url string) bool {
	regex, err := regexp.Compile("^(http|https)://")
	if err != nil {
		log.Println("Failed to compile URL validation regex:", err)
		return false
	}
	return regex.MatchString(strings.TrimSpace(url))
}
