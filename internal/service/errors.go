package service

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrTeamLimitReached  = errors.New("team limit reached")
	ErrForbidden         = errors.New("forbidden")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrCannotRemoveOwner = errors.New("cannot remove owner")
	ErrCannotArchive     = errors.New("cannot archive")
)
