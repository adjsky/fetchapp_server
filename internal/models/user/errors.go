package user

import "errors"

var (
	ErrEmailRegistered = errors.New("provided email is registered")
	ErrInternal        = errors.New("internal error")
	ErrInvalidToken    = errors.New("invalid auth token provided")
	ErrNotMatched      = errors.New("provided password doesn't match")
	ErrNoUser          = errors.New("no user with given email found")
)
