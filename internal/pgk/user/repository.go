package user

import (
	"gorm.io/gorm"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/repositories"
)

// repository interface
type Repository interface {
	Find(db *gorm.DB, req *Request) (*models.User, error)
	FindAll(db *gorm.DB, req *Request) ([]*models.User, error)
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
	if req.UserID != "" {
		db = db.Where("user_id = ?", req.UserID)
	}

	if req.Userlevel != "" {
		db = db.Where("user_level = ?", req.Userlevel)
	}

	if req.EmpID != "" {
		db = db.Where("emp_id = ?", req.EmpID)
	}

	if req.EmpID != "" {
		db = db.Where("emp_id = ?", req.EmpID)
	}

	if req.UserStatus != nil {
		db = db.Where("user_status IN ?", req.UserStatus)
	}

	return db
}

// Find find
func (r *repository) Find(db *gorm.DB, req *Request) (*models.User, error) {
	entities := &models.User{}
	query := r.query(db, req)
	err := query.Limit(1).Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}

// FindAll find all
func (r *repository) FindAll(db *gorm.DB, req *Request) ([]*models.User, error) {
	entities := []*models.User{}
	query := r.query(db, req)
	err := query.Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}
