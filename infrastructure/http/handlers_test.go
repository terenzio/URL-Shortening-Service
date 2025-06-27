package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/terenzio/URL-Shortening-Service/application"
	urlModel "github.com/terenzio/URL-Shortening-Service/domain"
)

// mockURLRepository is a simple mock for url repository used in tests.
type mockURLRepository struct {
	StoreFunc           func(ctx context.Context, url urlModel.URL) error
	FindByShortCodeFunc func(ctx context.Context, shortCode string) (*urlModel.URL, error)
	IsUniqueFunc        func(ctx context.Context, shortCode string) bool
	FetchAllFunc        func(ctx context.Context) ([]urlModel.URL, error)
}

func (m *mockURLRepository) Store(ctx context.Context, url urlModel.URL) error {
	if m.StoreFunc != nil {
		return m.StoreFunc(ctx, url)
	}
	return nil
}

func (m *mockURLRepository) FindByShortCode(ctx context.Context, shortCode string) (*urlModel.URL, error) {
	if m.FindByShortCodeFunc != nil {
		return m.FindByShortCodeFunc(ctx, shortCode)
	}
	return nil, errors.New("not implemented")
}

func (m *mockURLRepository) IsUnique(ctx context.Context, shortCode string) bool {
	if m.IsUniqueFunc != nil {
		return m.IsUniqueFunc(ctx, shortCode)
	}
	return true
}

func (m *mockURLRepository) FetchAll(ctx context.Context) ([]urlModel.URL, error) {
	if m.FetchAllFunc != nil {
		return m.FetchAllFunc(ctx)
	}
	return nil, nil
}

// helper to create test context
func newTestContext(method, path string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(method, path, bytes.NewReader(body))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	c.Request = req
	return c, w
}

func TestHandleHomePage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	expectedURL := urlModel.URL{ShortCode: "abc", OriginalURL: "https://example.com", Expiry: time.Now()}

	tests := []struct {
		name           string
		repo           *mockURLRepository
		expectedStatus int
		expectBody     bool
	}{
		{
			name: "success",
			repo: &mockURLRepository{
				FetchAllFunc: func(ctx context.Context) ([]urlModel.URL, error) {
					return []urlModel.URL{expectedURL}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name: "repo error",
			repo: &mockURLRepository{
				FetchAllFunc: func(ctx context.Context) ([]urlModel.URL, error) { return nil, errors.New("fail") },
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := application.NewURLService(tt.repo)
			h := NewHandler(service)

			c, w := newTestContext(http.MethodGet, "/url/display", nil)
			h.HandleHomePage(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectBody {
				var got []urlModel.URLMapping
				err := json.Unmarshal(w.Body.Bytes(), &got)
				assert.NoError(t, err)
				assert.Equal(t, expectedURL.ShortCode, got[0].ShortCode)
				assert.Equal(t, expectedURL.OriginalURL, got[0].OriginalURL)
			}
		})
	}
}

func TestHandleAddLink(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		body           []byte
		repoSetup      func(stored *urlModel.URL) *mockURLRepository
		expectedStatus int
		validate       func(t *testing.T, w *httptest.ResponseRecorder, stored urlModel.URL)
	}{
		{
			name: "bad request",
			body: []byte(`{}`),
			repoSetup: func(stored *urlModel.URL) *mockURLRepository {
				return &mockURLRepository{}
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid url",
			body: []byte(`{"original_url":"invalid"}`),
			repoSetup: func(stored *urlModel.URL) *mockURLRepository {
				return &mockURLRepository{}
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "custom short code not unique",
			body: []byte(`{"original_url":"https://example.com","custom_short_code":"dup"}`),
			repoSetup: func(stored *urlModel.URL) *mockURLRepository {
				return &mockURLRepository{
					IsUniqueFunc: func(ctx context.Context, code string) bool { return false },
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "custom short code success",
			body: []byte(`{"original_url":"https://example.com","custom_short_code":"mycode"}`),
			repoSetup: func(stored *urlModel.URL) *mockURLRepository {
				return &mockURLRepository{
					IsUniqueFunc: func(ctx context.Context, code string) bool { return true },
					StoreFunc: func(ctx context.Context, url urlModel.URL) error {
						*stored = url
						return nil
					},
				}
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder, stored urlModel.URL) {
				assert.Equal(t, "mycode", stored.ShortCode)
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Contains(t, resp["shortened_url"], "mycode")
			},
		},
		{
			name: "generate short code",
			body: []byte(`{"original_url":"https://example.com"}`),
			repoSetup: func(stored *urlModel.URL) *mockURLRepository {
				return &mockURLRepository{
					IsUniqueFunc: func(ctx context.Context, code string) bool { return true },
					StoreFunc: func(ctx context.Context, url urlModel.URL) error {
						*stored = url
						return nil
					},
				}
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder, stored urlModel.URL) {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Contains(t, resp["shortened_url"], stored.ShortCode)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stored urlModel.URL
			repo := tt.repoSetup(&stored)
			service := application.NewURLService(repo)
			h := NewHandler(service)

			c, w := newTestContext(http.MethodPost, "/url/add", tt.body)
			h.HandleAddLink(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validate != nil {
				tt.validate(t, w, stored)
			}
		})
	}
}

func TestHandleRedirectToOriginalLink(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		path           string
		shortcode      string
		repo           *mockURLRepository
		expectedStatus int
		expectedLoc    string
	}{
		{
			name:           "missing param",
			path:           "/redirect/",
			shortcode:      "",
			repo:           &mockURLRepository{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "not found",
			path:      "/redirect/abc",
			shortcode: "abc",
			repo: &mockURLRepository{
				FindByShortCodeFunc: func(ctx context.Context, code string) (*urlModel.URL, error) {
					return nil, errors.New("not found")
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:      "success",
			path:      "/redirect/abc",
			shortcode: "abc",
			repo: &mockURLRepository{
				FindByShortCodeFunc: func(ctx context.Context, code string) (*urlModel.URL, error) {
					return &urlModel.URL{OriginalURL: "https://example.com"}, nil
				},
			},
			expectedStatus: http.StatusTemporaryRedirect,
			expectedLoc:    "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := application.NewURLService(tt.repo)
			h := NewHandler(service)

			c, w := newTestContext(http.MethodGet, tt.path, nil)
			if tt.shortcode != "" {
				c.Params = gin.Params{{Key: "shortcode", Value: tt.shortcode}}
			}
			h.HandleRedirectToOriginalLink(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedLoc != "" {
				assert.Equal(t, tt.expectedLoc, w.Header().Get("Location"))
			}
		})
	}
}
