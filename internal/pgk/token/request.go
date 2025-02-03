package token

type Request struct {
	UserID    string `json:"user_id"`
	UserLevel string `json:"user_level"`
	EmpID     string `json:"emp_id"`
	SessionID string `json:"session_id"`
}

type RequestRefresh struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
