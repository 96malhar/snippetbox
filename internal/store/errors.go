package store

import "errors"

var (
	ErrNoRecord           = errors.New("store: no matching record found")
	ErrInvalidCredentials = errors.New("store: invalid credentials")
	ErrDuplicateEmail     = errors.New("store: duplicate email")
)
