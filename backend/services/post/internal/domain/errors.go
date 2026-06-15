package domain

import "errors"

var (
	ErrIDInvalid        = errors.New("invalid id")
	ErrContentBlank     = errors.New("content is blank")
	ErrPermissionDenied = errors.New("permission denied")
)
