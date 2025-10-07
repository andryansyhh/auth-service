package dto

import "errors"

var (
	ErrConflict     = errors.New("resource already exists")
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("invalid credentials")
	ErrInvalidInput = errors.New("invalid input provided")
)
