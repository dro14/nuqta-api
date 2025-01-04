package e

import "errors"

var (
	ErrNoParams = errors.New("params are not specified")
	ErrNotFound = errors.New("not found")
)
