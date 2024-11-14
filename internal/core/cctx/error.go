package cctx

import (
	"github.com/gofiber/fiber/v3"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
)

// NewError new custom fiber error
func (c *Context) NewError(code int, message string) *fiber.Error {
	return fiber.NewError(code, message)
}

// ErrorNotFound error not found `404`
func (c *Context) ErrorNotFound(message string) *fiber.Error {
	return c.NewError(config.RR.Internal.Unauthorized.HTTPStatusCode(), message)
}

// ErrorBadRequest error bad request `400`
func (c *Context) ErrorBadRequest(message string) *fiber.Error {
	return c.NewError(config.RR.Internal.Unauthorized.HTTPStatusCode(), message)
}

// ErrorUnauthorized error unauthorized `401`
func (c *Context) ErrorUnauthorized(message string) *fiber.Error {
	return c.NewError(config.RR.Internal.Unauthorized.HTTPStatusCode(), message)
}

// ErrorForbidden error forbidden `403`
func (c *Context) ErrorForbidden(message string) *fiber.Error {
	return c.NewError(config.RR.Internal.Unauthorized.HTTPStatusCode(), message)
}

// ErrorTooManyRequests error too many requests `429`
func (c *Context) ErrorTooManyRequests(message string) *fiber.Error {
	return c.NewError(config.RR.Internal.TooManyRequests.HTTPStatusCode(), message)
}
