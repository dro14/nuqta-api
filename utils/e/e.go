package e

import "errors"

var (
	ErrForbidden        = errors.New("forbidden")
	ErrInvalidMatch     = errors.New("invalid match")
	ErrInvalidParams    = errors.New("invalid params")
	ErrNoAuthHeader     = errors.New("no authorization header")
	ErrNoFilename       = errors.New("filename is not specified")
	ErrNoParams         = errors.New("params are not specified")
	ErrNotFound         = errors.New("not found")
	ErrUnknownEdge      = errors.New("unknown edge")
	ErrUnknownAttribute = errors.New("unknown attribute")
)

var (
	ErrBadRequest    = errors.New("bad request")
	ErrSpit          = errors.New("model spits")
	ErrStream        = errors.New("stream error")
	ErrEmpty         = errors.New("empty response")
	ErrDownload      = errors.New("download error")
	ErrTimeout       = errors.New("download timeout")
	ErrInappropriate = errors.New("inappropriate request")
	ErrContextLength = errors.New("context length exceeded")
)
