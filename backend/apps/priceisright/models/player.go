package models

type PlayerClass string

const (
	PLAYER_CLASS_AUDIENCE PlayerClass = "Audience"
	PLAYER_CLASS_DISPLAY  PlayerClass = "Display"
	PLAYER_CLASS_LEADER   PlayerClass = "Leader"
)

var PlayerClassStrings = []string{
	string(PLAYER_CLASS_LEADER),
	string(PLAYER_CLASS_DISPLAY),
	string(PLAYER_CLASS_AUDIENCE),
}

type Player struct {
	Class       PlayerClass
	Name        string
	BadgeNumber string
	UUID        string
}

func (P *Player) GetUUID() string {
	return P.UUID
}
