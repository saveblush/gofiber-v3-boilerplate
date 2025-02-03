package user

import (
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/cctx"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
)

// service interface
type Service interface {
	Find(c *cctx.Context, req *Request) (*models.User, error)
	FindAll(c *cctx.Context, req *Request) ([]*models.User, error)
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
func (s *service) Find(c *cctx.Context, req *Request) (*models.User, error) {
	res, err := s.repository.Find(c.GetRelayDatabase(), req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindAll find all
func (s *service) FindAll(c *cctx.Context, req *Request) ([]*models.User, error) {
	res, err := s.repository.FindAll(c.GetRelayDatabase(), req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
