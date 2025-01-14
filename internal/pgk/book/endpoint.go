package book

import (
	"github.com/gofiber/fiber/v3"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/handlers"
)

// endpoint interface
type Endpoint interface {
	Find(c fiber.Ctx) error
	FindByID(c fiber.Ctx) error
	Create(c fiber.Ctx) error
	Update(c fiber.Ctx) error
	Delete(c fiber.Ctx) error
}

type endpoint struct {
	config  *config.Configs
	result  *config.ReturnResult
	service Service
}

func NewEndpoint() Endpoint {
	return &endpoint{
		config:  config.CF,
		result:  config.RR,
		service: NewService(),
	}
}

// @Tags Book
// @Summary Get book
// @Description Get book
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Success 200 {object} models.Book
// @Failure 400 {object} models.Message
// @Failure 401 {object} models.Message
// @Failure 404 {object} models.Message
// @Failure 410 {object} models.Message
// @Security ApiKeyAuth
// @Router /book [get]
func (ep *endpoint) Find(c fiber.Ctx) error {
	return handlers.ResponseObject(c, ep.service.FindAll, &Request{})
}

// @Tags Book
// @Summary Get book
// @Description Get book
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Success 200 {object} models.Book
// @Failure 400 {object} models.Message
// @Failure 401 {object} models.Message
// @Failure 404 {object} models.Message
// @Failure 410 {object} models.Message
// @Security ApiKeyAuth
// @Router /book/{id} [get]
func (ep *endpoint) FindByID(c fiber.Ctx) error {
	return handlers.ResponseObject(c, ep.service.FindByID, &RequestID{})
}

// @Tags Book
// @Summary Create book
// @Description Create book
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Success 200 {object} models.Book
// @Failure 400 {object} models.Message
// @Failure 401 {object} models.Message
// @Failure 404 {object} models.Message
// @Failure 410 {object} models.Message
// @Security ApiKeyAuth
// @Router /book [post]
func (ep *endpoint) Create(c fiber.Ctx) error {
	return handlers.ResponseObject(c, ep.service.Create, &RequestCreate{})
}

// @Tags Book
// @Summary Update book
// @Description Update book
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Success 200 {object} models.Book
// @Failure 400 {object} models.Message
// @Failure 401 {object} models.Message
// @Failure 404 {object} models.Message
// @Failure 410 {object} models.Message
// @Security ApiKeyAuth
// @Router /book/:id [put]
func (ep *endpoint) Update(c fiber.Ctx) error {
	return handlers.ResponseObject(c, ep.service.Update, &RequestUpdate{})
}

// @Tags Book
// @Summary Delete book
// @Description Delete book
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Success 200 {object} models.Book
// @Failure 400 {object} models.Message
// @Failure 401 {object} models.Message
// @Failure 404 {object} models.Message
// @Failure 410 {object} models.Message
// @Security ApiKeyAuth
// @Router /book/:id [delete]
func (ep *endpoint) Delete(c fiber.Ctx) error {
	return handlers.ResponseSuccess(c, ep.service.Delete, &RequestID{})
}
