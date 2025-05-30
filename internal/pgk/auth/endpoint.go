package auth

import (
	"github.com/gofiber/fiber/v3"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/handlers"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/pgk/token"
)

// endpoint interface
type Endpoint interface {
	Login(c fiber.Ctx) error
	Logout(c fiber.Ctx) error
	Me(c fiber.Ctx) error
	RefreshToken(c fiber.Ctx) error
}

type endpoint struct {
	config  *config.Configs
	result  *config.ReturnResult
	service Service
}

func NewEndpoint() Endpoint {
	return &endpoint{
		config:  config.CF,
		result:  config.RR,
		service: NewService(),
	}
}

// @Tags Authentication
// @Summary Login
// @Description Login
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param Request body RequestLogin true "Body for login"
// @Success 200 {object} models.UserLogin
// @Failure 400 {object} models.Message
// @Failure 401 {object} models.Message
// @Failure 404 {object} models.Message
// @Failure 410 {object} models.Message
// @Security ApiKeyAuth
// @Router /login [post]
func (ep *endpoint) Login(c fiber.Ctx) error {
	return handlers.ResponseObject(c, ep.service.Login, &RequestLogin{})
}

// @Tags Authentication
// @Summary Logout
// @Description Logout
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Message
// @Failure 401 {object} models.Message
// @Failure 404 {object} models.Message
// @Failure 410 {object} models.Message
// @Security ApiKeyAuth
// @Router /logout [post]
func (ep *endpoint) Logout(c fiber.Ctx) error {
	return handlers.ResponseSuccessWithoutRequest(c, ep.service.Logout)
}

// @Tags Authentication
// @Summary Me
// @Description Me
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Message
// @Failure 401 {object} models.Message
// @Failure 404 {object} models.Message
// @Failure 410 {object} models.Message
// @Security ApiKeyAuth
// @Router /me [get]
func (ep *endpoint) Me(c fiber.Ctx) error {
	return handlers.ResponseSuccessWithoutRequest(c, ep.service.Me)
}

// @Tags Authentication
// @Summary RefreshToken
// @Description RefreshToken
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param request body token.RequestRefresh true "RefreshToken"
// @Success 200 {object} models.Token
// @Failure 400 {object} models.Message
// @Failure 401 {object} models.Message
// @Failure 404 {object} models.Message
// @Failure 410 {object} models.Message
// @Security ApiKeyAuth
// @Router /auth/refreshtoken [post]
func (ep *endpoint) RefreshToken(c fiber.Ctx) error {
	return handlers.ResponseObject(c, ep.service.RefreshToken, &token.RequestRefresh{})
}
