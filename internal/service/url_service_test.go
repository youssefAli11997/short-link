package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"url-shortener/internal/model"
	"url-shortener/internal/repository"
)

func TestURLService_Encode(t *testing.T) {
	dbError := errors.New("db down")

	tests := []struct {
		name         string
		originalURL  string
		createResult *model.URL
		createErr    error
		expectedURL  string
		expectedErr  error
	}{
		{
			name:        "invalid url",
			originalURL: "invalid",
			expectedErr: model.ErrInvalidURL,
		},
		{
			name:        "new url",
			originalURL: "https://google.com",
			createResult: &model.URL{
				ID:          100,
				OriginalURL: "https://google.com",
			},
			expectedURL: "http://localhost:8080/1C",
		},
		{
			name:        "url already exists",
			originalURL: "https://google.com",
			createResult: &model.URL{
				ID:          10,
				OriginalURL: "https://google.com",
			},
			expectedURL: "http://localhost:8080/a",
		},
		{
			name:        "repository fails",
			originalURL: "https://google.com",
			createErr:   dbError,
			expectedErr: dbError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &repository.MockURLRepository{
				CreateOrGetFunc: func(ctx context.Context, originalURL string) (*model.URL, error) {

					require.Equal(t, tt.originalURL, originalURL)

					return tt.createResult, tt.createErr
				},
			}

			service := NewURLService(mockRepo, "http://localhost:8080")

			got, err := service.Encode(context.Background(), tt.originalURL)

			require.ErrorIs(t, err, tt.expectedErr)
			require.Equal(t, tt.expectedURL, got)
		})
	}
}

func TestURLService_Decode(t *testing.T) {
	tests := []struct {
		name               string
		shortURL           string
		findResult         *model.URL
		findErr            error
		expectedID         int64
		expectedURL        string
		expectedErr        error
		shouldCallFindByID bool
	}{
		{
			name:        "invalid short url",
			shortURL:    "%%%%",
			expectedErr: model.ErrInvalidURL,
		},
		{
			name:        "invalid base62 code",
			shortURL:    "http://localhost:8080/$$$",
			expectedErr: model.ErrInvalidURL,
		},
		{
			name:               "not found",
			shortURL:           "http://localhost:8080/a",
			expectedID:         10,
			findErr:            model.ErrURLNotFound,
			expectedErr:        model.ErrURLNotFound,
			shouldCallFindByID: true,
		},
		{
			name:       "success single digit code",
			shortURL:   "http://localhost:8080/a",
			expectedID: 10,
			findResult: &model.URL{
				ID:          10,
				OriginalURL: "https://google.com",
			},
			expectedURL:        "https://google.com",
			shouldCallFindByID: true,
		},
		{
			name:       "success multi digit code",
			shortURL:   "http://localhost:8080/1C",
			expectedID: 100,
			findResult: &model.URL{
				ID:          100,
				OriginalURL: "https://github.com",
			},
			expectedURL:        "https://github.com",
			shouldCallFindByID: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var findByIDCalled bool

			mockRepo := &repository.MockURLRepository{
				FindByIDFunc: func(ctx context.Context, id int64) (*model.URL, error) {
					findByIDCalled = true

					require.Equal(t, tt.expectedID, id)

					return tt.findResult, tt.findErr
				},
			}

			service := NewURLService(mockRepo, "http://localhost:8080")

			got, err := service.Decode(context.Background(), tt.shortURL)

			require.ErrorIs(t, err, tt.expectedErr)
			require.Equal(t, tt.expectedURL, got)
			require.Equal(t, tt.shouldCallFindByID, findByIDCalled)
		})
	}
}
