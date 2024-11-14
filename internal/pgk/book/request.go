package book

import "github.com/saveblush/gofiber-v3-boilerplate/internal/request"

type Request struct {
	Name string `json:"name" query:"name"`
}

type RequestCreate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Author      string `json:"author"`
}

type RequestUpdate struct {
	request.GetOne
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Author      string `json:"author"`
}
