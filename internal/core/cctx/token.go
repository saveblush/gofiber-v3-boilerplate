package cctx

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/connection/cache"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/generic"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
)

// getClaims get user claims and session login
func (c *Context) getClaims() (*models.TokenClaims, *models.Token, error) {
	token, ok := c.Locals(UserKey).(*jwt.Token)
	if !ok {
		return nil, nil, errors.New("token claims not found")
	}

	// claim
	tokenClaims := token.Claims.(*models.TokenClaims)

	// get session login
	session := &models.Token{}
	if !generic.IsEmpty(tokenClaims.Subject) && !generic.IsEmpty(tokenClaims.SessionID) {
		ch := cache.New()
		key := utils.SetKeySessionLogin(tokenClaims.Subject, tokenClaims.SessionID)
		_ = ch.Get(key, session)
	}

	return tokenClaims, session, nil
}

// GetClaims get user claims
func (c *Context) GetClaims() (*models.TokenClaims, error) {
	token, _, err := c.getClaims()
	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetUserID get user id claims
func (c *Context) GetUserID() string {
	token, _, err := c.getClaims()
	if err != nil {
		return ""
	}

	return token.Subject
}

// GetUserLevel get user level claims
func (c *Context) GetUserLevel() string {
	_, session, err := c.getClaims()
	if err != nil {
		return ""
	}

	return session.User.UserLevel
}

// GetEmpID get emp id claims
func (c *Context) GetEmpID() string {
	_, session, err := c.getClaims()
	if err != nil {
		return ""
	}

	return session.User.EmpID
}
