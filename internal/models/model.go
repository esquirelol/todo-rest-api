package models

import "time"

type ModelTodo struct {
	Id          int
	Author      string
	Title       string
	Description string
	Status      bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}
