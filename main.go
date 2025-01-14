package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"

	"github.com/saveblush/gofiber-v3-boilerplate/docs"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/breaker"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/connection/sql"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/handlers/routes"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	flag.Parse()

	// Init logger
	logger.InitLogger()

	// Init configuration
	err := config.InitConfig()
	if err != nil {
		logger.Log.Fatalf("init configuration error: %s", err)
	}

	// Init return result
	err = config.InitReturnResult()
	if err != nil {
		logger.Log.Fatalf("init return result error: %s", err)
	}

	// Set swagger info
	docs.SwaggerInfo.Title = config.CF.Swagger.Title
	docs.SwaggerInfo.Description = config.CF.Swagger.Description
	docs.SwaggerInfo.Version = config.CF.Swagger.Version
	docs.SwaggerInfo.Host = fmt.Sprintf("%s%s", config.CF.Swagger.Host, config.CF.Swagger.BaseURL)
	//docs.SwaggerInfo.Schemes = []string{"https", "http"}

	// Init connection database
	if config.CF.Database.RelaySQL.Enable {
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
			logger.Log.Fatalf("init connection db error: %s", err)
		}
		sql.RelayDatabase = session.Database

		if !fiber.IsChild() {
			session.Database.AutoMigrate(&models.Book{})
		}
	}

	// Debug db
	if !config.CF.App.Environment.Production() {
		if config.CF.Database.RelaySQL.Enable {
			sql.DebugRelayDatabase()
		}
	}

	// Init Circuit Breaker
	breaker.Init()

	// Start app
	app := routes.NewServer()
	addr := flag.String("addr", fmt.Sprintf(":%d", config.CF.App.Port), "http service address")
	listenConfig := fiber.ListenConfig{
		EnablePrefork: config.CF.HTTPServer.Prefork,
	}
	go func() {
		err = app.Listen(*addr, listenConfig)
		if err != nil {
			logger.Log.Panicf("App start error: %s", err)
		}
	}()

	// Shutdown app
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Log.Info("Gracefully shutting down...")
	_ = app.Shutdown()

	// Cleanup tasks
	logger.Log.Info("Running cleanup tasks...")

	// Close db
	if config.CF.Database.RelaySQL.Enable {
		go sql.CloseConnection(sql.RelayDatabase)
	}
	logger.Log.Info("Database connection closed")

	logger.Log.Info("App was successful shutdown")
}
