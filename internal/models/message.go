package models

type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewSuccessMessage new success message
func NewSuccessMessage() *Message {
	return &Message{
		Code:    200,
		Message: "success",
	}
}
