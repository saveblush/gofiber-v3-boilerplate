package middlewares

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/gofiber/fiber/v3"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/cctx"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
)

// Logger logger
func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := utils.Now()
		err := c.Next()
		if err != nil {
			return err
		}

		userID := cctx.New(c).GetUserID()

		var b []byte
		parameters := c.Locals(cctx.ParametersKey)
		if parameters != nil {
			b, _ = utils.Marshal(parameters)
			for _, f := range []string{"Password", "password"} {
				if res := gjson.GetBytes(b, f); res.Exists() {
					b, _ = sjson.SetBytes(b, f, "**********")
				}
			}
		}

		logs := logger.Log().With(
			zap.String("host", c.Hostname()),
			zap.String("method", c.Method()),
			zap.String("path", c.OriginalURL()),
			zap.Any("language", c.Locals(cctx.LangKey)),
			zap.String("ip", c.IP()),
			zap.Any("ips", c.IPs()),
			zap.String("user_agent", c.Get(fiber.HeaderUserAgent)),
			zap.String("body_size", fmt.Sprintf("%.5f MB", float64(bytes.NewReader(c.Request().Body()).Len())/1024/1024)),
			zap.Any("process_time", time.Since(start)),
			zap.String("user_id", userID),
			zap.String("parameters", string(b)),
		)

		if c.OriginalURL() != fmt.Sprintf("%s/healthcheck", config.CF.App.ApiBaseUrl) {
			if !strings.HasPrefix(c.OriginalURL(), fmt.Sprintf("%s/swagger", config.CF.Swagger.BaseURL)) {
				logs.Infof("[%s][%s] response: %v", c.Method(), c.OriginalURL(), string(c.Response().Body()))
			}
		}

		return nil
	}
}
