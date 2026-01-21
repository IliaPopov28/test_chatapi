package validator

import (
	"errors"
)

const (
	MaxTitleLength = 200
	MaxTextLength  = 5000
	MinTitleLength = 1
	MinTextLength  = 1
	DefaultLimit   = 20
	MaxLimit       = 100
)

var (
	ErrTitleEmpty   = errors.New("title cannot be empty")
	ErrTitleTooLong = errors.New("title cannot be longer than 200 characters")
	ErrTextEmpty    = errors.New("text cannot be empty")
	ErrTextTooLong  = errors.New("text cannot be longer than 5000 characters")
	ErrInvalidLimit = errors.New("limit must be between 1 and 100")
)

func ValidateTitle(title string) error {
	if title == "" {
		return ErrTitleEmpty
	}
	if len(title) < MinTitleLength || len(title) > MaxTitleLength {
		return ErrTitleTooLong
	}
	return nil
}

func ValidateText(text string) error {
	if text == "" {
		return ErrTextEmpty
	}
	if len(text) < MinTextLength || len(text) > MaxTextLength {
		return ErrTextTooLong
	}
	return nil
}

func ValidateLimit(limit int) int {
	if limit <= 0 {
		return DefaultLimit
	}
	if limit > MaxLimit {
		return MaxLimit
	}
	return limit
}
