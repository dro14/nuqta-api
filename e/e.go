package e

import "errors"

var (
	ErrInvalidMatch     = errors.New("invalid match")
	ErrInvalidParams    = errors.New("invalid params")
	ErrNoAuthHeader     = errors.New("no authorization header")
	ErrNoFilename       = errors.New("filename is not specified")
	ErrNoParams         = errors.New("params are not specified")
	ErrNotFound         = errors.New("not found")
	ErrUnknownEdge      = errors.New("unknown edge")
	ErrUnknownPredicate = errors.New("unknown predicate")
)
