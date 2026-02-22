package routes

import (
	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/handlers/middlewares"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/pgk/auth"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/pgk/book"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/pgk/system"

	_ "github.com/saveblush/gofiber-v3-boilerplate/docs"
)

// InitRouter init router
func (s *server) InitRouter() {
	// api
	api := s.Group(s.config.App.ApiBaseUrl)

	// system
	systemEndpoint := system.NewEndpoint()
	systemRoute := api.Group("/system")
	systemRoute.Post("/action", middlewares.AuthorizationAdminRequired(), systemEndpoint.Action)
	systemRoute.Get("/maintenance", middlewares.Maintenance())

	api.Use(
		middlewares.Available(), // ปิด/เปิด ระบบ
		middlewares.AcceptLanguage(),
	)

	// healthcheck endpoint
	api.Get("/healthcheck", healthcheck.New())

	// api v1
	v1 := api.Group("/v1")

	// swagger
	if s.config.Swagger.Enable {
		v1.Get("/swagger/*", swaggo.New(swaggo.Config{
			Title: s.config.Swagger.Title,
		}))
	}

	// auth
	authEndpoint := auth.NewEndpoint()
	v1.Post("/login", middlewares.AuthorizationRequired(), authEndpoint.Login)
	v1.Post("/logout", middlewares.AuthorizationRequired(), authEndpoint.Logout)
	v1.Get("/me", middlewares.AuthorizationRequired(), authEndpoint.Me)

	authApi := v1.Group("auth")
	authApi.Use(middlewares.AuthorizationRequired())
	authApi.Post("/refreshtoken", authEndpoint.RefreshToken)

	// book
	bookEndpoint := book.NewEndpoint()
	bookApi := v1.Group("books", middlewares.AuthorizationRequired())
	bookApi.Get("", bookEndpoint.Find)
	bookApi.Get("/:id", bookEndpoint.FindByID)
	bookApi.Post("/", bookEndpoint.Create)
	bookApi.Put("/:id", bookEndpoint.Update)
	bookApi.Delete("/:id", bookEndpoint.Delete)

	// not found
	s.Use(middlewares.Notfound())
}
