package storage

import "errors"

var (
	ErrHashExists   = errors.New("hash already exists")
	ErrHashNotFound = errors.New("hash not found")
)
