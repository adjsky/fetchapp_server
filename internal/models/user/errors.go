package user

import "errors"

var (
	ErrEmailRegistered = errors.New("the provided email is registered")
	ErrInternal        = errors.New("internal error")
	ErrInvalidToken    = errors.New("an invalid auth token provided")
	ErrNotMatched      = errors.New("the provided password doesn't match the account password")
	ErrNoUser          = errors.New("no user with the given email found")
)
