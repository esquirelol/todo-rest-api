package models

import "time"

type ModelTodo struct {
	Author      string
	Title       string
	Description string
	Status      bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}
