package e

import "errors"

var (
	ErrNoParam      = errors.New("param is not specified")
	ErrNoQuery      = errors.New("query is not specified")
	ErrNoFilename   = errors.New("filename is not specified")
	ErrUnknownParam = errors.New("unknown param")
	ErrNotFound     = errors.New("not found")
)
