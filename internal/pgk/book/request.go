package book

import "github.com/saveblush/gofiber-v3-boilerplate/internal/models"

type RequestID struct {
	ID uint `json:"-" path:"id" form:"id" query:"id" validate:"required"`
}

type Request struct {
	ID         uint   `json:"id" query:"id"`
	IDs        []uint `json:"ids" query:"ids"`
	Name       string `json:"name" query:"name"`
	Pagination bool   `json:"pagination" query:"pagination"`
	models.PageForm
}

type RequestCreate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Author      string `json:"author"`
}

type RequestUpdate struct {
	RequestID
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Author      string `json:"author"`
}

type RequestAttach struct {
	RequestID
}
