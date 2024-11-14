package system

import (
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/cctx"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
)

// service interface
type Service interface {
	Action(c *cctx.Context, req *Request) (interface{}, error)
}

type service struct {
	config *config.Configs
	result *config.ReturnResult
}

func NewService() Service {
	return &service{
		config: config.CF,
		result: config.RR,
	}
}

// Action action
func (s *service) Action(c *cctx.Context, req *Request) (interface{}, error) {
	status := ""
	if req.Status == "on" {
		status = config.AvailableStatusOnline
	} else if req.Status == "off" {
		status = config.AvailableStatusOffline
	}

	if status == "" {
		return nil, s.result.Internal.BadRequest
	}

	// set config
	s.config.SetConfigAvailableStatus(status)
	s.config.SetConfigAvailableDescription(req.Body)

	return map[string]interface{}{
		"system": status,
	}, nil
}
