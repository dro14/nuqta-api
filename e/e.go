package e

import "errors"

var (
	ErrNoAuthHeader     = errors.New("no authorization header")
	ErrNoFilename       = errors.New("filename is not specified")
	ErrNoParams         = errors.New("params are not specified")
	ErrNotFound         = errors.New("not found")
	ErrUnknownEdge      = errors.New("unknown edge")
	ErrUnknownParam     = errors.New("unknown param")
	ErrUnknownPredicate = errors.New("unknown predicate")
)
