package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not exit")
	ErrURLExists   = errors.New("url exists")
)
