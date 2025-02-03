package auth

type RequestLogin struct {
	UserID   string `json:"user_id" validate:"required" example:"user"`
	Password string `json:"password" validate:"required" example:"password"`
}

type RequestBypass struct {
	Token string `json:"token" validate:"required" example:"sxos904ksdldsfigh4x849358gfj"`
}

type RequestLoginBypass struct {
	UserID string `json:"user_id" validate:"required" example:"user"`
}

type RequestLogLogin struct {
	SeqNo         string `json:"seq_no"`
	UserID        string `json:"user_id"`
	UserLevel     string `json:"user_level"`
	EmpID         string `json:"emp_id"`
	CompID        string `json:"comp_id"`
	ConnectIP     string `json:"connect_ip"`
	ConnectDevice string `json:"connect_device"`
	ConnectType   string `json:"connect_type"`
	ConnectResult string `json:"connect_result"`
}

type RequestLastLogin struct {
	UserID        string `json:"user_id"`
	ConnectIP     string `json:"connect_ip"`
	ConnectDevice string `json:"connect_device"`
}
