package user

type Request struct {
	Name string `json:"name" query:"name"`
}
