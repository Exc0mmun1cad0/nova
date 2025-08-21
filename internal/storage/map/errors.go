package mapstorage

import "errors"

var (
	ErrWrongType   = errors.New("wrong type operation")
	ErrKeyNotFound = errors.New("key not found")
)
