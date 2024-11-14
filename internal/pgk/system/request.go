package system

type Request struct {
	Status string `json:"status" validate:"required"`
	Body   string `json:"body" validate:"required"`
}
