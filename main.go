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
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/connection/cache"
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
		logger.Log.Panicf("init configuration error: %s", err)
	}

	// Init return result
	err = config.InitReturnResult()
	if err != nil {
		logger.Log.Panicf("init return result error: %s", err)
	}

	// Set swagger info
	docs.SwaggerInfo.Title = config.CF.Swagger.Title
	docs.SwaggerInfo.Description = config.CF.Swagger.Description
	docs.SwaggerInfo.Version = config.CF.Swagger.Version
	docs.SwaggerInfo.Host = fmt.Sprintf("%s%s", config.CF.Swagger.Host, config.CF.Swagger.BaseURL)
	//docs.SwaggerInfo.Schemes = []string{"https", "http"}

	// Init database
	err = initDatabase()
	if err != nil {
		logger.Log.Panicf("init database error: %s", err)
	}

	// Init cache
	err = initCache()
	if err != nil {
		logger.Log.Panicf("init cache error: %s", err)
	}

	// Init Circuit Breaker
	breaker.Init()

	// New app
	app, err := routes.NewServer()
	if err != nil {
		logger.Log.Panicf("new server error: %s", err)
	}

	// Start cron
	app.CronStart()

	// Init router
	app.InitRouter()

	// Listen app
	addr := flag.String("addr", fmt.Sprintf(":%d", config.CF.App.Port), "http service address")
	listenConfig := fiber.ListenConfig{
		EnablePrefork: config.CF.HTTPServer.Prefork,
	}
	go func() {
		err = app.Listen(*addr, listenConfig)
		if err != nil {
			logger.Log.Panicf("server start error: %s", err)
		}
	}()

	// Shutdown app
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	logger.Log.Info("Gracefully shutting down...")

	// Close app
	_ = app.Close()
	logger.Log.Info("Server closed")
	logger.Log.Info("Running cleanup tasks...")

	// Close cron
	app.CronStop()
	logger.Log.Info("Cron stoped")

	// Close cache
	_ = cache.New().Close()
	logger.Log.Info("Cache connection closed")

	// Close db
	_ = closeDatabase()
	logger.Log.Info("Database connection closed")

	logger.Log.Info("App was successful shutdown")
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
	sql.Database = session.Database

	if !fiber.IsChild() {
		sql.Database.AutoMigrate(
			&models.User{},
			&models.AuthLogLogin{},
			&models.Book{},
		)
	}

	// Debug db
	if !config.CF.App.Environment.Production() {
		sql.DebugDatabase()
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

// closeDatabase close connection database
func closeDatabase() error {
	sql.CloseConnection(sql.Database)

	return nil
}
