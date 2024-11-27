package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v3"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/handlers/render"
)

// Available available
// ปิด/เปิด ระบบ
func Available() fiber.Handler {
	return func(c fiber.Ctx) error {
		if config.CF.App.AvailableStatus == config.AvailableStatusOnline {
			return c.Next()
		} else {
			return fiber.NewError(fiber.StatusServiceUnavailable)
		}
	}
}

// Maintenance maintenance
// กรณีปิดระบบ ดึง body html
func Maintenance() fiber.Handler {
	return func(c fiber.Ctx) error {
		path := fmt.Sprintf("./templates/%s", config.CF.HTMLTemplate.SystemMaintenance)
		body, err := config.CF.ReadConfigAvailableDescription()
		if err != nil {
			logger.Log.Error("read file config available description error:", err)
			return fiber.NewError(fiber.StatusServiceUnavailable, "Error: Available Description")
		}

		return render.Html(c, path, body)
	}
}
