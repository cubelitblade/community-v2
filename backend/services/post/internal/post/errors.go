package post

import "errors"

var (
	ErrIDInvalid    = errors.New("invalid id")
	ErrContentBlank = errors.New("content is blank")
)
