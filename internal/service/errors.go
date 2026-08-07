package service

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrTeamLimitReached  = errors.New("team limit reached")
	ErrForbidden         = errors.New("forbidden")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrCannotRemoveOwner = errors.New("cannot remove owner")
	ErrCannotArchive     = errors.New("cannot archive")
	ErrAlreadyMember     = errors.New("alreaduy is a member")
	ErrAlreadyExists     = errors.New("already exists")
	ErrInvitationExpired = errors.New("invitation expired")
	ErrInvalidInviteLink = errors.New("invalid invite link")
	ErrAlreadyInvited    = errors.New("already invited")
)
