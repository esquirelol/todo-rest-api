package dto

import "time"

type TodoUpdate struct {
	Author      *string `json:"author,omitempty"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *bool   `json:"status,omitempty" `
}

type Todo struct {
	Author      string     `json:"author"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      bool       `json:"status"`
	CreatedAt   time.Time  `json:"created"`
	CompletedAt *time.Time `json:"completed"`
}
