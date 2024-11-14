package models

import "github.com/golang-jwt/jwt/v5"

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// token claims
type TokenClaims struct {
	jwt.RegisteredClaims
	Role  string `json:"role"`
	EmpID string `json:"emp_id"`
	LogId string `json:"log_id"` // id log connect
}

// ข้อมูล return สำหรับนำไป token.request ต่อ
type TokenUser struct {
	UserID    string `json:"user_id"`
	UserLevel string `json:"user_level"`
	EmpID     string `json:"emp_id"`
	LogId     string `json:"log_id"` // id log connect
}
