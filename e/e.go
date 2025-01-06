package e

import "errors"

var (
	ErrNoParams     = errors.New("params are not specified")
	ErrNoFilename   = errors.New("filename is not specified")
	ErrUnknownParam = errors.New("unknown param")
	ErrUnknownEdge  = errors.New("unknown edge")
	ErrNotFound     = errors.New("not found")
)
