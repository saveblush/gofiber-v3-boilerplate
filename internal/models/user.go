package models

import "time"

type User struct {
	UserID          string    `json:"user_id" gorm:"primaryKey;type:varchar(30)"`
	Password        string    `json:"-" gorm:"type:varchar(100)"`
	UserLevel       string    `json:"user_level" gorm:"type:varchar(30)"`
	Name            string    `json:"name"`
	Surname         string    `json:"surname"`
	Email           string    `json:"email"`
	EmpID           string    `json:"emp_id" gorm:"type:varchar(30)"`
	LastLoginAt     time.Time `json:"last_login_at" gorm:"index;type:timestamp without time zone"`
	LastLoginIP     string    `json:"last_login_ip"`
	LastLoginDevice string    `json:"last_login_device"`
	UserStatus      int       `json:"user_status" gorm:"type:smallint;default:1"`
}

func (User) TableName() string {
	return "al_com_user"
}

type UserLogin struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	SessionID    string `json:"session_id"`
	User         *User  `json:"user"`
	Code         int    `json:"code"`
	Message      string `json:"message"`
}
