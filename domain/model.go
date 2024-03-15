package domain

import "time"

// URL represents the URL entity in the domain layer
type URL struct {
	OriginalURL string    `json:"original_url"`
	Expiry      time.Time `json:"expiry"`
	ShortCode   string    `json:"short_code"`
}

// AddURLRequest represents the request body for adding a new URL.
type AddURLRequest struct {
	OriginalURL     string    `json:"original_url"`
	Expiry          time.Time `json:"expiry"`
	CustomShortCode string    `json:"custom_short_code"`
}

// AddSuccessResponse represents the response body for a successful URL addition.
type AddSuccessResponse struct {
	OriginalURL  string    `json:"original_url"`
	Expiry       time.Time `json:"expiry"`
	ShortenedURL string    `json:"shortened_url"`
}

// URLMapping represents the URL mapping entity in the domain layer.
// This is used to display the list of all shortened URLs.
type URLMapping struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	Expiry      time.Time `json:"expiry"`
}
