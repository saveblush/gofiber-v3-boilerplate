package token

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/cctx"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/generic"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
)

// Service service interface
type Service interface {
	Create(c *cctx.Context, req *Request) (*models.Token, error)
	VerifyRefresh(c *cctx.Context, token string) (*models.TokenUser, error)
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

// create token
func (s *service) Create(c *cctx.Context, req *Request) (*models.Token, error) {
	accessToken, err := s.genToken(req)
	if err != nil {
		logger.Log.Errorf("create access token error: %s", err)
		return nil, err
	}

	refreshToken, err := s.genRefreshToken(req)
	if err != nil {
		logger.Log.Errorf("create refresh token error: %s", err)
		return nil, err
	}

	return &models.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// verify refresh token
func (s *service) VerifyRefresh(c *cctx.Context, tokenString string) (*models.TokenUser, error) {
	if generic.IsEmpty(tokenString) {
		return nil, errors.New("refresh token not found")
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JWT.RefreshSecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*models.TokenClaims)
	if !ok {
		return nil, err
	}

	return &models.TokenUser{
		UserID:    claims.Subject,
		SessionID: claims.SessionID,
	}, nil
}

// gen token claim
func (s *service) genClaim(req *Request) *models.TokenClaims {
	claims := &models.TokenClaims{}
	claims.Issuer = s.config.App.Issuer
	claims.Subject = req.UserID
	claims.SessionID = req.SessionID

	return claims
}

// create jwt token
func (s *service) genToken(req *Request) (string, error) {
	now := utils.Now()
	claims := s.genClaim(req)
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(s.config.JWT.AccessExpireTime))

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(s.config.JWT.AccessSecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

// create jwt refresh token
func (s *service) genRefreshToken(req *Request) (string, error) {
	now := utils.Now()
	claims := s.genClaim(req)
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(s.config.JWT.RefreshExpireTime))

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(s.config.JWT.RefreshSecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}
