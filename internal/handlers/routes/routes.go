package routes

import (
	swagger "github.com/saveblush/gofiber3-swagger"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/handlers/middlewares"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/pgk/auth"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/pgk/book"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/pgk/healthcheck"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/pgk/system"
)

// InitRouter init router
func (s *server) InitRouter() {
	// api
	api := s.Group(s.config.App.ApiBaseUrl)

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
	if s.config.Swagger.Enable {
		v1.Get("/swagger/*", swagger.HandlerDefault)
	}

	// auth
	authEndpoint := auth.NewEndpoint()
	v1.Post("/login", authEndpoint.Login, middlewares.AuthorizationRequired())
	v1.Post("/logout", authEndpoint.Logout, middlewares.AuthorizationRequired())
	v1.Get("/me", authEndpoint.Me, middlewares.AuthorizationRequired())

	authApi := v1.Group("auth")
	authApi.Use(middlewares.AuthorizationRequired())
	authApi.Post("/refreshtoken", authEndpoint.RefreshToken)

	// book
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
