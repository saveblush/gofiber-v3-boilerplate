package auth

import (
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/cctx"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/pgk/token"
)

// Service service interface
type Service interface {
	Login(c *cctx.Context, req *Request) (interface{}, error)
}

type service struct {
	config *config.Configs
	result *config.ReturnResult
	token  token.Service
}

func NewService() Service {
	return &service{
		config: config.CF,
		result: config.RR,
		token:  token.NewService(),
	}
}

func (s *service) Login(c *cctx.Context, req *Request) (interface{}, error) {
	// authen

	// create token
	token, err := s.token.Create(c, &token.Request{
		// UserID:    dataUser.UserID,
		// UserLevel: dataUser.UserLevel,
		// EmpID:     dataUser.EmpID,
		// LogId:     logSeqno,
	})
	if err != nil {
		logger.Log.Errorf("create token error: %s", err)
		return nil, err
	}

	// เก็บ cache session login
	//s.createCacheSessionLogin(dataUser.UserID, logSeqno, &models.Token{AccessToken: token.AccessToken, RefreshToken: token.RefreshToken})

	res := &models.UserLogin{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return res, nil
}
