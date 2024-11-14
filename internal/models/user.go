package models

type User struct {
	Pubkey       string `json:"pubkey" gorm:"primaryKey;type:varchar(64)"`
	Name         string `json:"name"`
	LightningURL string `json:"lightning_url"`
}

func (User) TableName() string {
	return "users"
}

type UserLogin struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	LogConnectId string `json:"log_connect_id"`
	Code         int    `json:"code"`
	Message      string `json:"message"`
	User         *User  `json:"user"`
}
