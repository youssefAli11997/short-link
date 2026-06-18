package repository

import (
	"context"

	"url-shortener/internal/model"
)

type MockURLRepository struct {
	FindByOriginalURLFunc func(ctx context.Context, url string) (*model.URL, error)
	FindByShortCodeFunc   func(ctx context.Context, code string) (*model.URL, error)
	CreateFunc            func(ctx context.Context, originalURL string) (int64, error)
	UpdateShortCodeFunc   func(ctx context.Context, id int64, code string) error
}

func (m *MockURLRepository) FindByOriginalURL(ctx context.Context, url string) (*model.URL, error) {
	return m.FindByOriginalURLFunc(ctx, url)
}

func (m *MockURLRepository) FindByShortCode(ctx context.Context, code string) (*model.URL, error) {
	return m.FindByShortCodeFunc(ctx, code)
}

func (m *MockURLRepository) Create(ctx context.Context, originalURL string) (int64, error) {
	return m.CreateFunc(ctx, originalURL)
}

func (m *MockURLRepository) UpdateShortCode(ctx context.Context, id int64, code string) error {
	return m.UpdateShortCodeFunc(ctx, id, code)
}
