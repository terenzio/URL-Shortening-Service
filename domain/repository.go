package domain

import "context"

type URLRepository interface {
	Store(ctx context.Context, url URL) error
	FindByShortCode(ctx context.Context, shortCode string) (*URL, error)
	IsUnique(ctx context.Context, shortCode string) bool
	FetchAll(ctx context.Context) ([]URL, error)
}
