package routes

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/idempotency"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/handlers/middlewares"
)

const (
	// MaximumSize10MB body limit 1 mb.
	MaximumSize10MB = 10 * 1024 * 1024
	// MaximumSize1MB body limit 1 mb.
	MaximumSize1MB = 1 * 1024 * 1024
	// Timeout timeout 5 seconds
	Timeout5s = 5 * time.Second
	// Timeout timeout 10 seconds
	Timeout10s = 10 * time.Second
)

// NewServer new server
func NewServer() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:           config.CF.App.ProjectName,
		ServerHeader:      config.CF.App.ProjectName,
		BodyLimit:         MaximumSize10MB,
		ReadBufferSize:    MaximumSize1MB,
		WriteBufferSize:   MaximumSize1MB,
		IdleTimeout:       Timeout10s,
		ReadTimeout:       Timeout5s,
		WriteTimeout:      Timeout5s,
		ReduceMemoryUsage: true,
		CaseSensitive:     true,
		JSONEncoder:       utils.Marshal,
		JSONDecoder:       utils.Unmarshal,
	})

	// Middlewares
	app.Use(
		compress.New(compress.Config{
			Level: compress.LevelBestCompression,
		}),
		cors.New(),
		requestid.New(),
		idempotency.New(),
		pprof.New(),
		recover.New(),
	)

	// Logging
	/*app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))*/

	// Limiter
	if config.CF.HTTPServer.RateLimit.Enable {
		app.Use(limiter.New(limiter.Config{
			Max:        config.CF.HTTPServer.RateLimit.Max,
			Expiration: config.CF.HTTPServer.RateLimit.Expiration,
		}))
	}

	// Middlewares custom
	app.Use(
		middlewares.Logger(),
		middlewares.WrapError(),
	)

	// Setup the router
	NewRouter(app)

	return app
}
