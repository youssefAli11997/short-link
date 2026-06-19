package repository

import (
	"context"
	"url-shortener/internal/model"
)

type URLRepository interface {
	CreateOrGet(ctx context.Context, originalURL string) (*model.URL, error)
	FindByID(ctx context.Context, id int64) (*model.URL, error)
}
