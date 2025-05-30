package book

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/jinzhu/copier"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/cctx"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/connection/cache"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/generic"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
)

var (
	keyCache = "user"
)

// service interface
type Service interface {
	Find(c *cctx.Context, req *Request) (interface{}, error)
	FindAll(c *cctx.Context, req *Request) (interface{}, error)
	FindAllPage(c *cctx.Context, req *RequestPage) (interface{}, error)
	FindByID(c *cctx.Context, req *RequestID) (interface{}, error)
	Create(c *cctx.Context, req *RequestCreate) (interface{}, error)
	Update(c *cctx.Context, req *RequestUpdate) (interface{}, error)
	Delete(c *cctx.Context, req *RequestID) error
	Script(c *cctx.Context) error
}

type service struct {
	config     *config.Configs
	repository Repository
	cache      cache.Service
}

func NewService() Service {
	return &service{
		config:     config.CF,
		repository: NewRepository(),
		cache:      cache.New(),
	}
}

// Find find
func (s *service) Find(c *cctx.Context, req *Request) (interface{}, error) {
	res, err := s.repository.Find(c.GetDatabase(), req)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return map[string]interface{}{}, nil
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindAll find all
func (s *service) FindAll(c *cctx.Context, req *Request) (interface{}, error) {
	res, err := s.repository.FindAll(c.GetDatabase(), req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindAll find all page
func (s *service) FindAllPage(c *cctx.Context, req *RequestPage) (interface{}, error) {
	res, err := s.repository.FindAllPage(c.GetDatabase(), req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindByID find by id
func (s *service) FindByID(c *cctx.Context, req *RequestID) (interface{}, error) {
	key := fmt.Sprintf("%s-%d", keyCache, req.ID)
	res := &models.Book{}
	err := s.cache.Get(key, res)

	// ถ้าไม่เจอ cache
	if err != nil {
		err := s.repository.FindByID(c.GetDatabase(), req.ID, res)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]interface{}{}, nil
		}
		if err != nil {
			return nil, err
		}

		// เก็บใน cache
		if !generic.IsEmpty(req) {
			_ = s.cache.Set(key, res, s.config.Cache.ExprieTime.Default)
		}
	}

	return res, nil
}

// Create create
func (s *service) Create(c *cctx.Context, req *RequestCreate) (interface{}, error) {
	data := &models.Book{
		Name: req.Name,
	}

	err := s.repository.Create(c.GetDatabase(), data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Update update
func (s *service) Update(c *cctx.Context, req *RequestUpdate) (interface{}, error) {
	data := &models.Book{}
	copier.Copy(data, &req)

	err := s.repository.Update(c.GetDatabase(), data, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Delete delete
func (s *service) Delete(c *cctx.Context, req *RequestID) error {
	data := &models.Book{}
	copier.Copy(data, &req)

	err := s.repository.Delete(c.GetDatabase(), data)
	if err != nil {
		return err
	}

	return nil
}

// test cronjob
func (s *service) Script(c *cctx.Context) error {
	req := &RequestUpdate{
		RequestID: RequestID{ID: 2},
		Name:      "อัพเดทจาก cronjob",
	}
	data := &models.Book{}
	copier.Copy(data, req)

	err := s.repository.Update(c.GetDatabase(), data, data)
	if err != nil {
		return err
	}

	return nil
}
