package storage

import "errors"

var (
	ErrNoDataFound = errors.New("no data found")
	ErrUnexpected  = errors.New("unexpected error")
)
