package rest

import "github.com/pkg/errors"

var (
	ErrPostNotFound  = errors.New("post not found")
	ErrInvalidPostID = errors.New("invalid post id")
	ErrUnexpected    = errors.New("unexpected error")
)
