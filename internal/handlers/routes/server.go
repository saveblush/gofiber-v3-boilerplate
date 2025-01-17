package routes

import (
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/idempotency"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/robfig/cron/v3"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/breaker"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/cctx"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/connection/cache"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/connection/sql"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/handlers/middlewares"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
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
	var err error

	// Init database
	err = initDatabase()
	if err != nil {
		return nil, err
	}

	// Init cache
	err = initCache()
	if err != nil {
		return nil, err
	}

	// Init Circuit Breaker
	breaker.Init()

	// New app
	app := newApp()

	return &server{
		App:    app,
		cctx:   &cctx.Context{},
		config: config.CF,
		cron:   cron.New(),
	}, nil
}

// newApp new fiber app
func newApp() *fiber.App {
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
		JSONEncoder:       sonic.Marshal,
		JSONDecoder:       sonic.Unmarshal,
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

	return app
}

// initDatabase init connection database
func initDatabase() error {
	configuration := &sql.Configuration{
		Host:         config.CF.Database.RelaySQL.Host,
		Port:         config.CF.Database.RelaySQL.Port,
		Username:     config.CF.Database.RelaySQL.Username,
		Password:     config.CF.Database.RelaySQL.Password,
		DatabaseName: config.CF.Database.RelaySQL.DatabaseName,
		DriverName:   config.CF.Database.RelaySQL.DriverName,
		Charset:      config.CF.Database.RelaySQL.Charset,
		MaxIdleConns: config.CF.Database.RelaySQL.MaxIdleConns,
		MaxOpenConns: config.CF.Database.RelaySQL.MaxOpenConns,
		MaxLifetime:  config.CF.Database.RelaySQL.MaxLifetime,
	}
	session, err := sql.InitConnection(configuration)
	if err != nil {
		return err
	}
	sql.RelayDatabase = session.Database

	if !fiber.IsChild() {
		session.Database.AutoMigrate(&models.Book{})
	}

	// Debug db
	if !config.CF.App.Environment.Production() {
		sql.DebugRelayDatabase()
	}

	return nil
}

// initCache init cache
func initCache() error {
	configuration := &cache.Configuration{
		Host:     config.CF.Cache.Redis.Host,
		Port:     config.CF.Cache.Redis.Port,
		Password: config.CF.Cache.Redis.Password,
		DB:       config.CF.Cache.Redis.DB,
	}
	err := cache.Init(configuration)
	if err != nil {
		return err
	}

	return nil
}

// Close close server
func (s *server) Close() error {
	// Shutdown server
	err := s.Shutdown()
	if err != nil {
		return err
	}
	logger.Log.Info("Server closed")

	logger.Log.Info("Running cleanup tasks...")

	// Close cron
	s.CronStop()
	logger.Log.Info("Cron stoped")

	// Close cache
	cache.New().Close()
	logger.Log.Info("Cache connection closed")

	// Close db
	sql.CloseConnection(sql.RelayDatabase)
	logger.Log.Info("Database connection closed")

	return nil
}
