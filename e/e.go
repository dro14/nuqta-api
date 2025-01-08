package e

import "errors"

var (
	ErrNoParams     = errors.New("params are not specified")
	ErrNoFilename   = errors.New("filename is not specified")
	ErrNoAuthHeader = errors.New("no authorization header")
	ErrUnknownParam = errors.New("unknown param")
	ErrUnknownEdge  = errors.New("unknown edge")
	ErrNotFound     = errors.New("not found")
)
