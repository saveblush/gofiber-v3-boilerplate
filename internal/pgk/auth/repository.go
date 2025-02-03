package auth

import (
	"gorm.io/gorm"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/repositories"
)

// repository interface
type Repository interface {
	CreateLogLogin(db *gorm.DB, req *RequestLogLogin) (string, error)
	UpdateLastLogin(db *gorm.DB, req *RequestLastLogin) error
}

type repository struct {
	repositories.Repository
}

func NewRepository() Repository {
	return &repository{
		repositories.NewRepository(),
	}
}

// CreateLogLogin create log login
func (r *repository) CreateLogLogin(db *gorm.DB, req *RequestLogLogin) (string, error) {
	if req.CompID == "0" {
		req.CompID = "" // ถ้าเป็นค่า "0" จะ set เป็นค่าว่าง (เพราะแปลงจาก int)
	}

	data := &models.AuthLogLogin{
		SeqNo:         utils.UUID(),
		UserID:        req.UserID,
		UserLevel:     req.UserLevel,
		EmpID:         req.EmpID,
		CompID:        req.CompID,
		ConnectAt:     utils.Now(),
		ConnectIP:     req.ConnectIP,
		ConnectDevice: req.ConnectDevice,
		ConnectType:   req.ConnectType,
		ConnectResult: req.ConnectResult,
		Status:        1,
	}
	err := db.Create(data).Error
	if err != nil {
		return "", err
	}

	return data.SeqNo, nil
}

// UpdateLastLogin update last login
func (r *repository) UpdateLastLogin(db *gorm.DB, req *RequestLastLogin) error {
	data := &models.User{
		LastLoginAt:     utils.Now(),
		LastLoginIP:     req.ConnectIP,
		LastLoginDevice: req.ConnectDevice,
	}
	err := db.Model(&models.User{}).Where("user_id = ?", req.UserID).Updates(data).Error
	if err != nil {
		return err
	}

	return nil
}
