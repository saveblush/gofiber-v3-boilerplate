package auth

type Request struct {
	UserName string `json:"username" validate:"required" example:"10056789"`
	Password string `json:"password" validate:"required" example:"P@ssword1234"`
	Mode     string `json:"-"`
}
