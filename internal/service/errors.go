package service

import "errors"

var (
	ErrUserNotExist      = errors.New("user not found")
	ErrInvalidBalance    = errors.New("invalid balance: balance negative")
	ErrInvalidNewBalance = errors.New("invalid new balance: balance negative")
)
