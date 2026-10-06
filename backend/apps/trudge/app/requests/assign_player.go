package requests

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
)

type AssignPlayerRequest struct {
	UUID        string
	PlayerClass models.PlayerClass
}
