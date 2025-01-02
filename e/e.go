package e

import "errors"

var (
	ErrNoId     = errors.New("id is not specified")
	ErrNotFound = errors.New("not found")
)
