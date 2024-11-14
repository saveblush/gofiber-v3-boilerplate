package user

import (
	"gorm.io/gorm"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/repositories"
)

// repository interface
type Repository interface {
	Find(db *gorm.DB, req *Request) error
	FindByIDString(db *gorm.DB, field string, value string, i interface{}) error
}

type repository struct {
	repositories.Repository
}

func NewRepository() Repository {
	return &repository{
		repositories.NewRepository(),
	}
}

// Find find
func (r *repository) Find(db *gorm.DB, req *Request) error {

	return nil
}
