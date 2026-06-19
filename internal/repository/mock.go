package repository

import (
	"context"

	"url-shortener/internal/model"
)

type MockURLRepository struct {
	CreateOrGetFunc func(ctx context.Context, originalURL string) (*model.URL, error)
	FindByIDFunc    func(ctx context.Context, id int64) (*model.URL, error)
}

func (m *MockURLRepository) CreateOrGet(ctx context.Context, originalURL string) (*model.URL, error) {
	return m.CreateOrGetFunc(ctx, originalURL)
}

func (m *MockURLRepository) FindByID(ctx context.Context, id int64) (*model.URL, error) {
	return m.FindByIDFunc(ctx, id)
}
