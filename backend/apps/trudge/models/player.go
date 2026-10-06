package models

type PlayerClass string

const (
	PLAYER_CLASS_AUDIENCE PlayerClass = "Audience"
	PLAYER_CLASS_BARD     PlayerClass = "Bard"
	PLAYER_CLASS_DISPLAY  PlayerClass = "Display"
	PLAYER_CLASS_GOD      PlayerClass = "God"
	PLAYER_CLASS_LEADER   PlayerClass = "Leader"
	PLAYER_CLASS_WARRIOR  PlayerClass = "Warrior"
	PLAYER_CLASS_WIZARD   PlayerClass = "Wizard"
)

var TrudgePlayerClassesStrings = []string{
	string(PLAYER_CLASS_LEADER),
	string(PLAYER_CLASS_DISPLAY),
	string(PLAYER_CLASS_GOD),
	string(PLAYER_CLASS_AUDIENCE),
	string(PLAYER_CLASS_WARRIOR),
	string(PLAYER_CLASS_BARD),
	string(PLAYER_CLASS_WIZARD),
}

type Player struct {
	CharacterName string
	CharacterRace string
	Class         PlayerClass
	Equipment     []string
	Experience    int
	Health        int
	Name          string
	Option        int
	RoomKey       string
	UUID          string
	Vote          int
}

func (P *Player) GetUUID() string {
	return P.UUID
}
