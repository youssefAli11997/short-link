package service

import (
	"context"
	"errors"
	"testing"
	"url-shortener/internal/model"
	"url-shortener/internal/repository"
)

// TODO(ya): consider using testify/require for simpler test functions
// here and in the other test files

func TestURLService_Encode(t *testing.T) {
	dbError := errors.New("db down")

	tests := []struct {
		name        string
		originalURL string
		findResult  *model.URL
		findErr     error
		createID    int64
		createErr   error
		updateErr   error
		expectedURL string
		expectedErr error
	}{
		{
			name:        "invalid url",
			originalURL: "invalid",
			expectedErr: model.ErrInvalidURL,
		},
		{
			name:        "url already exists",
			originalURL: "https://google.com",
			findResult: &model.URL{
				ID:          10,
				OriginalURL: "https://google.com",
				ShortCode:   "a",
			},
			expectedURL: "http://localhost:8080/a",
		},
		{
			name:        "new url",
			originalURL: "https://google.com",
			findErr:     model.ErrURLNotFound,
			createID:    100,
			expectedURL: "http://localhost:8080/1C",
		},
		{
			name:        "create fails",
			originalURL: "https://google.com",
			findErr:     model.ErrURLNotFound,
			createErr:   dbError,
			expectedErr: dbError,
		},
		{
			name:        "update fails",
			originalURL: "https://google.com",
			updateErr:   dbError,
			createID:    100,
			expectedErr: dbError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &repository.MockURLRepository{
				FindByOriginalURLFunc: func(ctx context.Context, url string) (*model.URL, error) {
					return tt.findResult, tt.findErr
				},
				CreateFunc: func(ctx context.Context, originalURL string) (int64, error) {
					return tt.createID, tt.createErr
				},
				UpdateShortCodeFunc: func(ctx context.Context, id int64, code string) error {
					if id != 100 {
						t.Errorf("expected id 100, got %d", id)
					}

					if code != "1C" {
						t.Errorf("expected code 1C, got %s", code)
					}
					return tt.updateErr
				},
			}

			service := NewURLService(mockRepo, "http://localhost:8080")

			got, err := service.Encode(context.Background(), tt.originalURL)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v got %v", tt.expectedErr, err)
			}

			if got != tt.expectedURL {
				t.Fatalf("expected %s got %s", tt.expectedURL, got)
			}
		})
	}
}

func TestURLService_Decode(t *testing.T) {
	tests := []struct {
		name        string
		shortURL    string
		findResult  *model.URL
		findErr     error
		expectedURL string
		expectedErr error
	}{
		{
			name:        "invalid short url",
			shortURL:    "%%%%",
			expectedErr: model.ErrInvalidURL,
		},
		{
			name:        "not found",
			shortURL:    "http://localhost:8080/a",
			findErr:     model.ErrURLNotFound,
			expectedErr: model.ErrURLNotFound,
		},
		{
			name:     "success",
			shortURL: "http://localhost:8080/a",
			findResult: &model.URL{
				OriginalURL: "https://google.com",
			},
			expectedURL: "https://google.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &repository.MockURLRepository{
				FindByShortCodeFunc: func(ctx context.Context, code string) (*model.URL, error) {
					return tt.findResult, tt.findErr
				},
			}

			service := NewURLService(mockRepo, "http://localhost:8080")

			got, err := service.Decode(context.Background(), tt.shortURL)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected %v got %v", tt.expectedErr, err)
			}

			if got != tt.expectedURL {
				t.Fatalf("expected %s got %s", tt.expectedURL, got)
			}
		})
	}
}
