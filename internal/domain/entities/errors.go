package domain

import "errors"

var (
	ErrEventNotFound          = errors.New("event not found")
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidUserID          = errors.New("invalid user ID")
	ErrInvalidUserName        = errors.New("invalid user name")
	ErrInvalidEventTitle      = errors.New("invalid event title")
	ErrRegistrationNotFound   = errors.New("registration not found")
	ErrParticipantNotFound    = errors.New("participant not found")
	ErrConcurrentModification = errors.New("concurrent modification detected")
	ErrAdminOnly              = errors.New("only admins can perform this action")
)
