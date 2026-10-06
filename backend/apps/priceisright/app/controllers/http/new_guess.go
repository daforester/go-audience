package http

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/processors"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NewGuessController struct {
	controller.HTTP
	processor processors.NewGuessProcessorInterface
}

func (C NewGuessController) New(processor processors.NewGuessProcessorInterface) *NewGuessController {
	c := new(NewGuessController)
	c.processor = processor

	return c
}

func (C *NewGuessController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.NewGuess(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *NewGuessController) NewGuess(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	newGuess := &struct {
		Guess int64
	}{
		0,
	}
	err = c.BindJSON(newGuess)
	if err != nil {
		return controller.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	guess, err := C.processor.NewGuess(player, newGuess.Guess)

	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Guess *models.Guess
	}{
		guess,
	})
}
