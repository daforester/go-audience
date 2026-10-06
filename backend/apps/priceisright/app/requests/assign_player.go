package requests

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
)

type AssignPlayerRequest struct {
	UUID        string
	PlayerClass models.PlayerClass
}
