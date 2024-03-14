package application

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"github.com/terenzio/URL-Shortening-Service/domain"
)

const base62Characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

type URLService struct {
	repo domain.URLRepository
}

// NewURLService creates a new instance of URLService.
func NewURLService(repo domain.URLRepository) *URLService {
	return &URLService{repo: repo}
}

// ShortenURL generates a unique short code for the given URL and stores it in the repository.
func (s *URLService) ShortenURL(ctx context.Context, originalURL string, expiryTime time.Time) (string, error) {
	sequence := 1
	var shortCode string
	for {
		// Generate a unique short code
		shortCode = generateShortCode(originalURL, sequence)
		// Check if the short code is unique
		if s.repo.IsUnique(ctx, shortCode) {
			break
		}
		sequence++
	}

	url := domain.URL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		Expiry:      expiryTime,
	}

	// Store the URL with the generated short code in the repository
	if err := s.repo.Store(ctx, url); err != nil {
		return "", fmt.Errorf("failed to store URL: %w", err)
	}

	return shortCode, nil
}

// generateShortCode creates a unique short code for the given URL.
// It uses a sequence number to handle hash collisions and ensure uniqueness.
// The function generates a SHA-256 hash of the URL and encodes it using Base62.
// The result is a short code of up to 8 characters, which is then checked for uniqueness.
// If the short code is not unique, the sequence number is incremented and the process is repeated.
func generateShortCode(url string, sequence int) string {
	hash := hashURL(url, sequence)
	partialHash := hash[:10]
	shortCode := base62Encode(partialHash)
	if len(shortCode) > 8 {
		shortCode = shortCode[:8]
	}
	return shortCode
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

	// Ensure the result is 8 characters by padding with "0" (or another character),
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

// FetchAllURLs retrieves all URLs from the repository.
func (s *URLService) FetchAllURLs(ctx context.Context) ([]domain.URL, error) {
	return s.repo.FetchAll(ctx)
}

// GetOriginalURL retrieves the original URL for the given short code from the repository.
func (s *URLService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	url, err := s.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return "", fmt.Errorf("failed to find URL by short code: %w", err)
	}
	return url.OriginalURL, nil
}
