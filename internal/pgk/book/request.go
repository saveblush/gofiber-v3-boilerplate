package book

type RequestID struct {
	ID uint `json:"-" path:"id" form:"id" query:"id" validate:"required"`
}

type Request struct {
	Name string `json:"name" query:"name"`
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
