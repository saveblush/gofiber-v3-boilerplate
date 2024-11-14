package book

import (
	"gorm.io/gorm"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/generic"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/repositories"
)

// repository interface
type Repository interface {
	Find(db *gorm.DB, req *Request) (*models.Book, error)
	FindByID(db *gorm.DB, id uint, i interface{}) error
	Create(db *gorm.DB, i interface{}) error
	Update(db *gorm.DB, m, i interface{}) error
	Delete(db *gorm.DB, i interface{}) error
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
func (r *repository) Find(db *gorm.DB, req *Request) (*models.Book, error) {
	entities := &models.Book{}
	query := db

	if generic.IsEmpty(req.Name) {
		query.Where("name = ?", req.Name)
	}

	err := query.First(entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}
