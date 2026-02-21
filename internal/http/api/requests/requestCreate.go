package requests

type RequestCreate struct {
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
