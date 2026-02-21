package todo

import "time"

type Todo struct {
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created"`
	CompletedAt time.Time `json:"completed"`
}
