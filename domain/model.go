package domain

import "time"

type URL struct {
	ShortCode   string
	OriginalURL string
	Expiry      time.Time
}
