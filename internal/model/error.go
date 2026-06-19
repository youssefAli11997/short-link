package model

import "errors"

var (
	ErrURLNotFound         = errors.New("url not found")
	ErrInvalidURL          = errors.New("invalid url")
	ErrInvalidRequestBody  = errors.New("invalid request body")
	ErrAlreadyExists       = errors.New("record already exists")
	ErrEmptyBase62String   = errors.New("empty base62 string")
	ErrInternalServerError = errors.New("internal server error")
)
