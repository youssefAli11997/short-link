package service

import (
	"context"
	"net/url"
	"strings"
	"url-shortener/internal/model"
	"url-shortener/internal/repository"
	"url-shortener/internal/shortener"
)

type URLService interface {
	Encode(ctx context.Context, originalURL string) (string, error)
	Decode(ctx context.Context, shortURL string) (string, error)
}

type urlService struct {
	repository repository.URLRepository
	baseURL    string
}

func NewURLService(repository repository.URLRepository, baseURL string) URLService {
	return &urlService{
		repository: repository,
		baseURL:    baseURL,
	}
}

// TODO(ya): put creating and updating the db record in one database transaction
// to avoid race conditions
func (s *urlService) Encode(ctx context.Context, originalURL string) (string, error) {
	// 1. validate the url
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return "", model.ErrInvalidURL
	}

	// 2. create a row in the database and get the id
	url, err := s.repository.CreateOrGet(ctx, originalURL)
	if err != nil {
		return "", err
	}

	// 3. generate short code
	shortCode := shortener.EncodeBase62(url.ID)

	// 4. return short url
	return s.constructShortUrl(shortCode), nil
}

func (s *urlService) Decode(ctx context.Context, shortURL string) (string, error) {
	// 1. validate the url
	parsedUrl, err := url.Parse(shortURL)
	if err != nil {
		return "", model.ErrInvalidURL
	}

	// 2. extract short code from the url
	shortCode := strings.TrimPrefix(parsedUrl.Path, "/")

	// 3. decode short code to get id
	id, err := shortener.DecodeBase62(shortCode)
	if err != nil {
		return "", model.ErrInvalidURL
	}

	// 4. look up in db
	record, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return "", model.ErrURLNotFound
	}

	// 5. return original url
	return record.OriginalURL, nil
}

func (s *urlService) constructShortUrl(shortCode string) string {
	return s.baseURL + "/" + shortCode
}
