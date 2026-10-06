package models

type PlayerClass string

const (
	PLAYER_CLASS_ADMIN    PlayerClass = "Admin"
	PLAYER_CLASS_AUDIENCE PlayerClass = "Audience"
	PLAYER_CLASS_DIRECTOR PlayerClass = "Director"
	PLAYER_CLASS_DISPLAY  PlayerClass = "Display"
	PLAYER_CLASS_STAFF    PlayerClass = "Staff"
)

var FantasyPlayerClassesStrings = []string{
	string(PLAYER_CLASS_ADMIN),
	string(PLAYER_CLASS_DISPLAY),
	string(PLAYER_CLASS_STAFF),
	string(PLAYER_CLASS_AUDIENCE),
	string(PLAYER_CLASS_DIRECTOR),
}

func ValidPlayerClass(class PlayerClass) bool {
	for _, c := range FantasyPlayerClassesStrings {
		if c == string(class) {
			return true
		}
	}
	return false
}

type Player struct {
	Class   PlayerClass
	Name    string
	RoomKey string
	TeamID  string
	UUID    string
}

func (P *Player) GetUUID() string {
	return P.UUID
}

func (P *Player) Copy() *Player {
	p := *P
	return &p
}
