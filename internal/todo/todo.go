package todo

import "time"

type Todo struct {
	Author       string    `json:"author"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Status       bool      `json:"status"`
	Created_at   time.Time `json:"created"`
	Completed_at time.Time `json:"completed"`
}
