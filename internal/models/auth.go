package models

import "time"

type AuthLogLogin struct {
	SeqNo         string    `json:"seq_no" gorm:"primaryKey;type:char(36)"`
	UserID        string    `json:"user_id" gorm:"index;type:varchar(30)"`
	UserLevel     string    `json:"user_level" gorm:"index;type:varchar(30)"`
	EmpID         string    `json:"emp_id" gorm:"index;type:varchar(30)"`
	CompID        string    `json:"comp_id" gorm:"index;type:varchar(30)"`
	ConnectAt     time.Time `json:"connect_at" gorm:"index;type:timestamp without time zone"`
	ConnectIP     string    `json:"connect_ip"`
	ConnectDevice string    `json:"connect_device"`
	ConnectType   string    `json:"connect_type" gorm:"type:varchar(30)"`
	ConnectResult string    `json:"connect_result" gorm:"type:varchar(30)"`
	Status        int       `json:"status" gorm:"type:smallint;default:1"`
}

func (AuthLogLogin) TableName() string {
	return "log_connect"
}
