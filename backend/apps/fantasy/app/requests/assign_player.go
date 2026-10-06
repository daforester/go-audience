package requests

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
)

type AssignPlayerRequest struct {
	UUID        string
	PlayerClass models.PlayerClass
	TeamID      string
}
