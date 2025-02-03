package user

type Request struct {
	UserID     string `json:"user_id"`
	Userlevel  string `json:"user_level"`
	EmpID      string `json:"emp_id"`
	UserStatus []int  `json:"user_status"`
}
