package routes

import (
	swagger "github.com/saveblush/gofiber3-swagger"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/handlers/middlewares"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/pgk/book"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/pgk/healthcheck"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/pgk/system"
)

// InitRouter init router
func (s *server) InitRouter() {
	// api
	api := s.Group(config.CF.App.ApiBaseUrl)

	// system
	systemEndpoint := system.NewEndpoint()
	systemRoute := api.Group("/system")
	systemRoute.Post("/action", systemEndpoint.Action, middlewares.AuthorizationAdminRequired())
	systemRoute.Get("/maintenance", middlewares.Maintenance())

	api.Use(
		middlewares.Available(), // ปิด/เปิด ระบบ
		middlewares.AcceptLanguage(),
	)

	// healthcheck endpoint
	healthCheckEndpoint := healthcheck.NewEndpoint()
	api.Get("/healthcheck", healthCheckEndpoint.HealthCheck)

	// api v1
	v1 := api.Group("/v1")

	// swagger
	if config.CF.Swagger.Enable {
		v1.Get("/swagger/*", swagger.HandlerDefault)
	}

	v1.Get("/healthcheck", healthCheckEndpoint.HealthCheck)

	bookEndpoint := book.NewEndpoint()
	bookApi := v1.Group("book", middlewares.AuthorizationRequired())
	bookApi.Get("", bookEndpoint.Find)
	bookApi.Get("/:id", bookEndpoint.FindByID)
	bookApi.Post("/", bookEndpoint.Create)
	bookApi.Put("/:id", bookEndpoint.Update)
	bookApi.Delete("/:id", bookEndpoint.Delete)

	// not found
	s.Use(middlewares.Notfound())
}
