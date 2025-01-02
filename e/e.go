package e

import "errors"

var (
	ErrNoUID    = errors.New("uid is not specified")
	ErrNotFound = errors.New("not found")
)
