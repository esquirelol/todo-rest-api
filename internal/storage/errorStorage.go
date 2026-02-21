package storage

import "errors"

var (
	ErrExists       = errors.New("task is exists")
	ErrTaskNotFound = errors.New("task is not found")
)
