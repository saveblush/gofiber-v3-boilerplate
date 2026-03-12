package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/timeout"
)

// Timeout timeout
func Timeout(handler fiber.Handler) fiber.Handler {
	return timeout.New(handler, timeout.Config{Timeout: 10 * time.Second})
}
