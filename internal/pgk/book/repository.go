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
	FindAllPage(db *gorm.DB, req *Request) (*models.Page, error)
	Create(db *gorm.DB, i interface{}) error
	Update(db *gorm.DB, m, i interface{}) error
	Delete(db *gorm.DB, i interface{}) error
	DeleteFile(db *gorm.DB, req *RequestAttach) error
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

	if !generic.IsEmpty(req.IDs) {
		db = db.Where("id IN ?", req.IDs)
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
	err := query.Limit(1).Order("id").
		Preload("Display", "attach_type = ?", "1").
		Find(entities).Error
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

// FindAll find all page
func (r *repository) FindAllPage(db *gorm.DB, req *Request) (*models.Page, error) {
	entities := []*models.Book{}
	query := r.query(db, req)
	page, err := r.FindAllAndPageInformation(query, &req.PageForm, &entities)
	if err != nil {
		return nil, err
	}

	return models.NewPage(page, entities), nil
}

func (r *repository) DeleteFile(db *gorm.DB, req *RequestAttach) error {
	err := r.SoftDelete(db, "book_id", req.ID, "", models.BookFiles{})
	if err != nil {
		return err
	}

	return nil
}
