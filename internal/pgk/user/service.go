package user

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/cctx"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/generic"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
)

// service interface
type Service interface {
	Find(c *cctx.Context, req *Request) (interface{}, error)
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
	resNotfound := map[string]interface{}{
		"status": "error",
	}

	if generic.IsEmpty(req.Name) {
		res := resNotfound
		res["message"] = "field validation for 'name'"

		return res, nil
	}

	fetch := &models.User{}
	err := s.repository.FindByIDString(c.GetRelayDatabase(), "name", req.Name, fetch)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var res map[string]interface{}
	if !generic.IsEmpty(fetch.Name) && !generic.IsEmpty(fetch.Pubkey) {
		res = map[string]interface{}{
			"names": map[string]interface{}{
				fetch.Name: fetch.Pubkey,
			},
		}
	} else {
		res = resNotfound
		res["message"] = fmt.Sprintf("%s is not found", req.Name)
	}

	return res, nil
}
