package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terenzio/URL-Shortening-Service/domain"

	"github.com/go-redis/redis/v8"
)

type URLRepository struct {
	client *redis.Client
}

func NewURLRepository(client *redis.Client) *URLRepository {
	return &URLRepository{client: client}
}

// Store saves a URL entity to Redis, setting an expiry based on the URL's Expiry field.
// The short code is used as the key to store the original URL.
func (r *URLRepository) Store(ctx context.Context, url domain.URL) error {
	// Calculate the TTL (time-to-live) for the Redis entry based on the URL's expiry.
	// If the expiry is in the past, return an error.
	ttl := url.Expiry.Sub(time.Now())
	if ttl <= 0 {
		return fmt.Errorf("invalid expiry for URL %s", url.OriginalURL)
	}

	// Use the short code as the key to store the original URL.
	status := r.client.Set(ctx, "short:"+url.ShortCode, url.OriginalURL, ttl)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

// FindByShortCode retrieves a URL by its short code from Redis.
func (r *URLRepository) FindByShortCode(ctx context.Context, shortCode string) (*domain.URL, error) {
	result, err := r.client.Get(ctx, "short:"+shortCode).Result()
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("short code not found: %s", shortCode)
	} else if err != nil {
		return nil, err
	}

	return &domain.URL{ShortCode: shortCode, OriginalURL: result}, nil
}

// IsUnique checks if a short code is unique by attempting to find it in Redis.
func (r *URLRepository) IsUnique(ctx context.Context, shortCode string) bool {
	exists, err := r.client.Exists(ctx, "short:"+shortCode).Result()
	if err != nil {
		log.Printf("Error checking uniqueness in Redis: %v", err)
		return false
	}
	return exists == 0 // 0 means the key does not exist in Redis (i.e., it is unique)
}

func (r *URLRepository) FetchAll(ctx context.Context) ([]domain.URL, error) {
	var urls []domain.URL

	// Implementation depends on how URLs are stored
	// This is a simplified example. Consider using SCAN for production use.
	iter := r.client.Scan(ctx, 0, "short:*", 0).Iterator()
	for iter.Next(ctx) {
		shortCode := iter.Val()
		originalURL, err := r.client.Get(ctx, shortCode).Result()
		if err != nil {
			return nil, err
		}
		shortCode = strings.TrimPrefix(shortCode, "short:")
		urls = append(urls, domain.URL{ShortCode: shortCode, OriginalURL: originalURL})
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}
