package storage

import "errors"

var (
	ErrNotExists    = errors.New("task is not exists")
	ErrTaskNotFound = errors.New("task is not found")
)
