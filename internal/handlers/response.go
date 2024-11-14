package handlers

import (
	"reflect"

	"github.com/gofiber/fiber/v3"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/cctx"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/handlers/render"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
)

// ResponseObject handle response object
func ResponseObject(c fiber.Ctx, fn interface{}, request interface{}) error {
	ctx := cctx.New(c)
	err := ctx.BindValue(request, true)
	if err != nil {
		logger.Log().Errorf("bind value error: %s", err)
		return err
	}

	out := reflect.ValueOf(fn).Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(request),
	})

	errObj := out[1].Interface()
	if errObj != nil {
		logger.Log().Errorf("call service error: %s", errObj)
		return errObj.(error)
	}

	return render.JSON(c, out[0].Interface())
}

// ResponseObjectWithoutRequest handle response object without request
func ResponseObjectWithoutRequest(c fiber.Ctx, fn interface{}) error {
	ctx := cctx.New(c)
	out := reflect.ValueOf(fn).Call([]reflect.Value{
		reflect.ValueOf(ctx),
	})

	errObj := out[1].Interface()
	if errObj != nil {
		logger.Log().Errorf("call service error: %s", errObj)
		return errObj.(error)
	}

	return render.JSON(c, out[0].Interface())
}

// ResponseSuccess handle response success
func ResponseSuccess(c fiber.Ctx, fn interface{}, request interface{}) error {
	ctx := cctx.New(c)
	err := ctx.BindValue(request, true)
	if err != nil {
		logger.Log().Errorf("bind value error: %s", err)
		return err
	}

	out := reflect.ValueOf(fn).Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(request),
	})

	errObj := out[0].Interface()
	if errObj != nil {
		logger.Log().Errorf("call service error: %s", errObj)
		return errObj.(error)
	}

	return render.JSON(c, models.NewSuccessMessage())
}

// ResponseSuccessWithoutRequest handle response success without request
func ResponseSuccessWithoutRequest(c fiber.Ctx, fn interface{}) error {
	ctx := cctx.New(c)
	out := reflect.ValueOf(fn).Call([]reflect.Value{
		reflect.ValueOf(ctx),
	})

	errObj := out[0].Interface()
	if errObj != nil {
		logger.Log().Errorf("call service error: %s", errObj)
		return errObj.(error)
	}

	return render.JSON(c, models.NewSuccessMessage())
}
