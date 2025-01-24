package routes

import (
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/idempotency"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/robfig/cron/v3"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/cctx"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/handlers/middlewares"
)

const (
	// MaximumSize10MB body limit 1 mb.
	MaximumSize10MB = 10 * 1024 * 1024
	// MaximumSize1MB body limit 1 mb.
	MaximumSize1MB = 1 * 1024 * 1024
	// Timeout timeout 120 seconds
	Timeout120s = 120 * time.Second
	// Timeout timeout 10 seconds
	Timeout10s = 10 * time.Second
)

type server struct {
	// fiber
	*fiber.App

	// core context
	cctx *cctx.Context

	// config
	config *config.Configs

	// cron
	cron *cron.Cron
}

// NewServer new server
func NewServer() (*server, error) {
	// New fiber app
	app := fiber.New(fiber.Config{
		AppName:           config.CF.App.ProjectName,
		ServerHeader:      config.CF.App.ProjectName,
		BodyLimit:         MaximumSize10MB,
		IdleTimeout:       Timeout120s,
		ReadTimeout:       Timeout10s,
		WriteTimeout:      Timeout10s,
		ReduceMemoryUsage: true,
		CaseSensitive:     true,
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
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

	return &server{
		App:    app,
		cctx:   &cctx.Context{},
		config: config.CF,
		cron:   cron.New(),
	}, nil
}

// Close close server
func (s *server) Close() error {
	// Shutdown server
	err := s.Shutdown()
	if err != nil {
		return err
	}

	return nil
}
