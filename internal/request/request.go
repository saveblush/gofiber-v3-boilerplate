package request

// GetOne get one
type GetOne struct {
	ID uint `json:"-" path:"id" form:"id" query:"id" validate:"required"`
}

// GetOneString get one string
type GetOneString struct {
	ID string `json:"-" path:"id" form:"id" query:"id" validate:"required"`
}
