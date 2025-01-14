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
	FindAll(db *gorm.DB, req *Request) ([]*models.Book, error)
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

func (r *repository) query(db *gorm.DB, req *Request) *gorm.DB {
	if !generic.IsEmpty(req.ID) {
		db = db.Where("id = ?", req.ID)
	}

	if !generic.IsEmpty(req.Name) {
		db = db.Where("name = ?", req.Name)
	}

	return db
}

// Find find
func (r *repository) Find(db *gorm.DB, req *Request) (*models.Book, error) {
	entities := &models.Book{}
	query := r.query(db, req)
	err := query.First(entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}

// FindAll find all
func (r *repository) FindAll(db *gorm.DB, req *Request) ([]*models.Book, error) {
	entities := []*models.Book{}
	query := r.query(db, req)
	err := query.Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}
