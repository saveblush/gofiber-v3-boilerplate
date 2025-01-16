package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"

	"github.com/saveblush/gofiber-v3-boilerplate/docs"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/handlers/routes"
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

	// Start app
	app, err := routes.NewServer()
	if err != nil {
		logger.Log.Panicf("new server error: %s", err)
	}
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
	_ = app.Close()
}
