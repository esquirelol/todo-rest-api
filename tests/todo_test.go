package tests

import (
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/esquirelol/todo-rest-api/internal/dto"
	"github.com/esquirelol/todo-rest-api/internal/models"
	"github.com/gavv/httpexpect/v2"
)

const (
	host = "localhost:9091"
)

func TestTodo_Create(t *testing.T) {
	u := url.URL{
		Host:   host,
		Scheme: "http",
	}

	e := httpexpect.Default(t, u.String())
	e.POST("/create").WithJSON(dto.Todo{
		Author:      gofakeit.Name(),
		Title:       gofakeit.Sentence(3),
		Description: gofakeit.Sentence(5),
	}).Expect().Status(200)
}

func TestTodo_Get(t *testing.T) {
	u := url.URL{
		Host:   host,
		Scheme: "http",
	}
	author := gofakeit.Name()
	e := httpexpect.Default(t, u.String())
	e.GET("/{author}").WithPath("author", author).Expect().Status(404)
}

func TestTodo_CreateUpdateDelete(t *testing.T) {
	var data models.ModelTodo
	u := url.URL{
		Host:   host,
		Scheme: "http",
	}
	e := httpexpect.Default(t, u.String())
	author := gofakeit.Name()
	e.POST("/create").WithJSON(dto.Todo{
		Author:      author,
		Title:       gofakeit.Sentence(3),
		Description: gofakeit.Sentence(5),
	}).Expect().Status(200)
	e.GET("/{author}").WithPath("author", author).Expect().Status(200).JSON().Array().Value(0).Object().Decode(&data)
	e.PATCH("/{id}").WithPath("id", data.Id).WithJSON(dto.Todo{
		Title:       gofakeit.Sentence(3),
		Description: gofakeit.Sentence(5),
		Status:      true,
	}).Expect()
	e.DELETE("/{id}").WithPath("id", data.Id).Expect().Status(200)
}
