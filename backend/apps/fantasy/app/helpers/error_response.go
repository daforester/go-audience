package helpers

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/processors"
	"bitbucket.org/daforester/go-audience/apps/fantasy/services/controller"
	"errors"
	"github.com/gofiber/fiber/v3"
)

// ErrorResponse maps processor errors onto HTTP status codes
func ErrorResponse(err error) *controller.Response {
	switch {
	case errors.Is(err, processors.ErrForbidden):
		return controller.ErrorResponse(fiber.StatusForbidden, err.Error())
	case errors.Is(err, processors.ErrGameNotFound), errors.Is(err, processors.ErrPlayerNotFound):
		return controller.ErrorResponse(fiber.StatusNotFound, err.Error())
	case errors.Is(err, processors.ErrWrongPhase):
		return controller.ErrorResponse(fiber.StatusConflict, err.Error())
	case errors.Is(err, processors.ErrTooSoon):
		return controller.ErrorResponse(fiber.StatusTooManyRequests, err.Error())
	case errors.Is(err, processors.ErrInvalid):
		return controller.ErrorResponse(fiber.StatusBadRequest, err.Error())
	}
	return controller.ErrorResponse(fiber.StatusInternalServerError, err.Error())
}
