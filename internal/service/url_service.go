package service

import (
	"context"
)

type URLService interface {
	Encode(ctx context.Context, originalURL string) (string, error)
	Decode(ctx context.Context, shortURL string) (string, error)
}

func NewURLService() URLService {
	return &MockURLService{}
}
