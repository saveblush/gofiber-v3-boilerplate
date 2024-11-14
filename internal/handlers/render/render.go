package render

import (
	"github.com/gofiber/fiber/v3"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
)

// JSON render json to client
func JSON(c fiber.Ctx, response interface{}) error {
	return c.
		Status(fiber.StatusOK).
		JSON(response)
}

// Byte render byte to client
func Byte(c fiber.Ctx, bytes []byte) error {
	_, err := c.Status(fiber.StatusOK).
		Write(bytes)

	return err
}

// Error render error to client
func Error(c fiber.Ctx, err error) error {
	if fiberErr, ok := err.(*fiber.Error); ok {
		errMsg := config.RR.CustomMessage(fiberErr.Error(), fiberErr.Error(), fiberErr.Code)
		return c.
			Status(errMsg.Code).
			JSON(errMsg.WithLocale(c))
	}

	errMsg := config.RR.Internal.ConnectionError
	if locErr, ok := err.(config.Result); ok {
		errMsg = locErr
	}

	return c.
		Status(errMsg.HTTPStatusCode()).
		JSON(errMsg.WithLocale(c))
}

// Html render html to client
func Html(c fiber.Ctx, path, body string) error {
	return c.Render(path, fiber.Map{
		"body": body,
	})
}
