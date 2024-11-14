package middlewares

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/handlers/render"
)

// HandlerError handler error fiber
var HandlerError = func(c fiber.Ctx, err error) error {
	var e *fiber.Error
	if errors.As(err, &e) {
		return render.Error(c, err)
	}

	return nil
}

// WrapError wrap error
func WrapError() fiber.Handler {
	return func(c fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			return render.Error(c, err)
		}

		return nil
	}
}

// Notfound not found route
func Notfound() fiber.Handler {
	return func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "Not Found")
	}
}
