package store

import "errors"

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrPretendToCandelize = errors.New("pretending to candelize data")
)
