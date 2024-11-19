package cctx

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
)

// GetClaims get user claims
func (c *Context) GetClaims() (*models.TokenClaims, error) {
	token, ok := c.Locals(UserKey).(*jwt.Token)
	if !ok {
		return nil, errors.New("token claims not found")
	}

	return token.Claims.(*models.TokenClaims), nil
}

// GetUserID get user id claims
func (c *Context) GetUserID() string {
	token, err := c.GetClaims()
	if err != nil {
		return ""
	}

	return token.Subject
}

// GetUserLevel get user level claims
func (c *Context) GetUserLevel() string {
	token, err := c.GetClaims()
	if err != nil {
		return ""
	}

	return token.Role
}
