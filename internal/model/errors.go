package model

import "errors"

var (
	ErrURLNotFound        = errors.New("url not found")
	ErrInvalidURL         = errors.New("invalid url")
	ErrInvalidRequestBody = errors.New("invalid request body")
)
