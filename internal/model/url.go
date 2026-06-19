package model

import (
	"time"
)

type URL struct {
	ID          int64
	OriginalURL string
	CreatedAt   time.Time
}
