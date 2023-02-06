package app

import "errors"

var (
	ErrBadRequest = errors.New("bad request")
	ErrIDNotFound = errors.New("id not found")
	ErrNotFound   = errors.New("not found")
)
