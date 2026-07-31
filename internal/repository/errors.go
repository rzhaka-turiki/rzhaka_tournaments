package repository

import "errors"

var (
	ErrNotFound  = errors.New("repository: entity not found")
	ErrForbidden = errors.New("forbidden")
)
