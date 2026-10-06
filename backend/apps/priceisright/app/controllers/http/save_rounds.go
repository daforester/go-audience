package http

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/helpers"
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/processors"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/services/controller"
	"bitbucket.org/daforester/go-audience/services/ginhelper"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type SaveRoundsController struct {
	controller.HTTP
	processor      processors.SaveRoundsProcessorInterface
	imageProcessor processors.SaveImageUploadsProcessorInterface
}

func (C SaveRoundsController) New(processor processors.SaveRoundsProcessorInterface, imageProcessor processors.SaveImageUploadsProcessorInterface) *SaveRoundsController {
	c := new(SaveRoundsController)
	c.processor = processor
	c.imageProcessor = imageProcessor

	return c
}

func (C *SaveRoundsController) Run() error {
	c := C.GetContext().(*gin.Context)
	e := C.SaveRounds(c)

	c.JSON(
		e.ResponseCode,
		e,
	)

	return nil
}

func (C *SaveRoundsController) SaveRounds(c *gin.Context) *controller.Response {
	player, err := helpers.GetRequestPlayer(c)

	if err != nil {
		return controller.ErrorResponse(http.StatusUnauthorized, err.Error())
	}

	if player.Class != models.PLAYER_CLASS_LEADER {
		return controller.ErrorResponse(http.StatusUnauthorized, "must be leader")
	}

	input, err := ginhelper.GetRequestBody(c)

	rounds := make([]*models.Round, len(input["UUID"].(map[string]interface{})))
	files := make([]*multipart.FileHeader, len(rounds))

	for i, round := range rounds {
		l := strconv.Itoa(i)
		round = new(models.Round)
		rounds[i] = round

		round.UUID = input["UUID"].(map[string]interface{})[l].(string)
		round.Name = input["Name"].(map[string]interface{})[l].(string)
		round.Value, _ = strconv.ParseInt(input["Value"].(map[string]interface{})[l].(string), 10, 64)

		f := input["imageFile"].(map[string]interface{})[l]
		switch f.(type) {
		case *multipart.FileHeader:
			files[i] = f.(*multipart.FileHeader)
		case multipart.FileHeader:
			fx := f.(multipart.FileHeader)
			files[i] = &fx
		}
	}

	path := viper.GetString("ROOT_PATH") + "/public/images/rounds/"
	err = os.MkdirAll(path, 0666)
	if err != nil {
		return controller.ErrorResponse(500, "couldn't create path "+path)
	}

	err = C.imageProcessor.SaveImageUploads(rounds, files, path)
	if err != nil {
		return controller.ErrorResponse(500, "couldn't save images: "+err.Error())
	}
	rounds, err = C.processor.SaveRounds(rounds)

	if err != nil {
		return controller.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return controller.StandardResponse(struct {
		Rounds []*models.Round
	}{
		rounds,
	})
}
