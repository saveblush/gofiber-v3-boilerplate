package book

import (
	"errors"

	"gorm.io/gorm"

	"github.com/jinzhu/copier"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/cctx"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
)

// service interface
type Service interface {
	Find(c *cctx.Context, req *Request) (interface{}, error)
	FindByID(c *cctx.Context, req *RequestID) (interface{}, error)
	Create(c *cctx.Context, req *RequestCreate) (interface{}, error)
	Update(c *cctx.Context, req *RequestUpdate) (interface{}, error)
	Delete(c *cctx.Context, req *RequestID) error
}

type service struct {
	config     *config.Configs
	repository Repository
}

func NewService() Service {
	return &service{
		config:     config.CF,
		repository: NewRepository(),
	}
}

// Find find
func (s *service) Find(c *cctx.Context, req *Request) (interface{}, error) {
	res, err := s.repository.Find(c.GetRelayDatabase(), req)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return map[string]interface{}{}, nil
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindByID find by id
func (s *service) FindByID(c *cctx.Context, req *RequestID) (interface{}, error) {
	res := &models.Book{}
	err := s.repository.FindByID(c.GetRelayDatabase(), req.ID, res)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return map[string]interface{}{}, nil
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

// Create create
func (s *service) Create(c *cctx.Context, req *RequestCreate) (interface{}, error) {
	data := &models.Book{
		Name: req.Name,
	}

	err := s.repository.Create(c.GetRelayDatabase(), data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Update update
func (s *service) Update(c *cctx.Context, req *RequestUpdate) (interface{}, error) {
	data := &models.Book{}
	copier.Copy(data, &req)

	err := s.repository.Update(c.GetRelayDatabase(), data, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Delete delete
func (s *service) Delete(c *cctx.Context, req *RequestID) error {
	data := &models.Book{}
	copier.Copy(data, &req)

	err := s.repository.Delete(c.GetRelayDatabase(), data)
	if err != nil {
		return err
	}

	return nil
}
