package posts

import "errors"

var (
	ErrInvalidQueryParam = errors.New("invalid query param")
	ErrUnexpected        = errors.New("unexpected error")
)
