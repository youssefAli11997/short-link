package service

import "context"

type MockURLService struct{}

func (m *MockURLService) Encode(ctx context.Context, originalURL string) (string, error) {
	return "http://localhost:8080/abc123", nil
}

func (m *MockURLService) Decode(ctx context.Context, shortURL string) (string, error) {
	return "https://google.com", nil
}
