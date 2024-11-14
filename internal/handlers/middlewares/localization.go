package middlewares

import (
	"github.com/gofiber/fiber/v3"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/cctx"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
)

// AcceptLanguage header Accept-Language
func AcceptLanguage() fiber.Handler {
	return func(c fiber.Ctx) error {
		value := c.Get(fiber.HeaderAcceptLanguage)
		lang := config.Language(value)
		if lang != config.LanguageEN && lang != config.LanguageTH {
			lang = config.LanguageTH
		}
		c.Locals(cctx.LangKey, lang)

		return c.Next()
	}
}
