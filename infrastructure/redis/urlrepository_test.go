package redis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/terenzio/URL-Shortening-Service/domain"
)

func TestURLRepository_Store(t *testing.T) {
	// Setup a mini Redis server
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' occurred when starting miniredis", err)
	}
	defer mr.Close()

	// Connect to mini Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	repo := NewURLRepository(rdb)

	tests := []struct {
		name         string
		shortCode    string
		originalURL  string
		expiryOffset time.Duration // Offset from now
		wantErr      bool
	}{
		{
			name:         "Valid URL with 24-hour Expiry",
			shortCode:    "abc123",
			originalURL:  "https://example.com",
			expiryOffset: 24 * time.Hour,
			wantErr:      false,
		},
		{
			name:         "Expired URL",
			shortCode:    "expired123",
			originalURL:  "https://expired.com",
			expiryOffset: -1 * time.Hour, // 1 hour in the past
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := domain.URL{
				ShortCode:   tt.shortCode,
				OriginalURL: tt.originalURL,
				Expiry:      time.Now().Add(tt.expiryOffset),
			}

			// Attempt to store the URL
			err := repo.Store(context.Background(), url)
			if tt.wantErr {
				assert.Error(t, err, "Expected an error for test case: %s", tt.name)
			} else {
				assert.NoError(t, err, "Unexpected error for test case: %s", tt.name)

				// Verify data is stored only if no error is expected
				stored, err := mr.Get("short:" + tt.shortCode)
				if !tt.wantErr { // Check only if we don't expect an error
					assert.NoError(t, err, "Retrieving data from miniredis should not produce an error for test case: %s", tt.name)
					assert.Equal(t, tt.originalURL, stored, "Stored URL should match for test case: %s", tt.name)
				}
			}
		})
	}
}

func TestURLRepository_FindByShortCode(t *testing.T) {
	// Setup a mini Redis server
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' occurred when starting miniredis", err)
	}
	defer mr.Close()

	// Connect to mini Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	repo := NewURLRepository(rdb)

	// Prepopulate Redis with a test URL
	shortCode := "abc123"
	originalURL := "https://example.com"
	expiry := 24 * time.Hour
	mr.Set("short:"+shortCode, originalURL)
	mr.SetTTL("short:"+shortCode, expiry)

	tests := []struct {
		name      string
		shortCode string
		wantURL   *domain.URL
		wantErr   bool
	}{
		{
			name:      "URL Found",
			shortCode: shortCode,
			wantURL: &domain.URL{
				ShortCode:   shortCode,
				OriginalURL: originalURL,
				Expiry:      time.Now().Add(expiry), // Mocked expiry might not match exactly
			},
			wantErr: false,
		},
		{
			name:      "URL Not Found",
			shortCode: "nonExistent",
			wantURL:   nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, err := repo.FindByShortCode(context.Background(), tt.shortCode)
			if tt.wantErr {
				assert.Error(t, err, "FindByShortCode should return an error for test case: %s", tt.name)
			} else {
				assert.NoError(t, err, "FindByShortCode should not return an error for test case: %s", tt.name)
				assert.Equal(t, tt.wantURL.ShortCode, gotURL.ShortCode, "ShortCode should match for test case: %s", tt.name)
				assert.Equal(t, tt.wantURL.OriginalURL, gotURL.OriginalURL, "OriginalURL should match for test case: %s", tt.name)
				// Note: Direct comparison of Expiry might not be reliable due to precision differences.
			}
		})
	}
}

func TestURLRepository_IsUnique(t *testing.T) {
	// Setup a mini Redis server
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' occurred when starting miniredis", err)
	}
	defer mr.Close()

	// Connect to mini Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	repo := NewURLRepository(rdb)

	// Prepopulate Redis with a test URL to simulate a non-unique scenario
	existingShortCode := "existing123"
	mr.Set("short:"+existingShortCode, "https://example.com")

	tests := []struct {
		name       string
		shortCode  string
		wantUnique bool
	}{
		{
			name:       "ShortCode is Unique",
			shortCode:  "unique123",
			wantUnique: true,
		},
		{
			name:       "ShortCode is Not Unique",
			shortCode:  existingShortCode,
			wantUnique: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isUnique := repo.IsUnique(context.Background(), tt.shortCode)
			assert.Equal(t, tt.wantUnique, isUnique, "IsUnique should correctly determine uniqueness for test case: %s", tt.name)
		})
	}
}

func TestURLRepository_FetchAll(t *testing.T) {
	// Setup a mini Redis server
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' occurred when starting miniredis", err)
	}
	defer mr.Close()

	// Connect to mini Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	repo := NewURLRepository(rdb)

	// Prepopulate Redis with test URLs
	testURLs := []domain.URL{
		{ShortCode: "code1", OriginalURL: "https://example1.com", Expiry: time.Now().Add(24 * time.Hour)},
		{ShortCode: "code2", OriginalURL: "https://example2.com", Expiry: time.Now().Add(24 * time.Hour)},
	}
	for _, url := range testURLs {
		mr.Set("short:"+url.ShortCode, url.OriginalURL)
		// Note: Expiry is not directly stored in this example; adjust according to your schema.
	}

	tests := []struct {
		name     string
		wantURLs []domain.URL
		wantErr  bool
	}{
		{
			name:     "Successfully Fetch All URLs",
			wantURLs: testURLs,
			wantErr:  false,
		},
		// Additional test scenarios could include testing with an empty database or simulating a database error.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURLs, err := repo.FetchAll(context.Background())
			if tt.wantErr {
				assert.Error(t, err, "FetchAll should return an error for test case: %s", tt.name)
			} else {
				assert.NoError(t, err, "FetchAll should not return an error for test case: %s", tt.name)
				// Compare the lengths of the slices as a basic check
				assert.Equal(t, len(tt.wantURLs), len(gotURLs), "The number of fetched URLs should match for test case: %s", tt.name)
				// More detailed comparisons can be added as needed, depending on how much control you have over ordering and what fields are available
			}
		})
	}
}
