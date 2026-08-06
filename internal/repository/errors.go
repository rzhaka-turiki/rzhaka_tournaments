package repository

import "errors"

var (
	ErrNotFound      = errors.New("repository: entity not found")
	ErrConflict      = errors.New("conflict")
	ErrInvalid       = errors.New("invalid")
	ErrNotTeamMember = errors.New("user is not team member")
)
