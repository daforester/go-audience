package models

import "time"

type Guess struct {
	PlayerUUID string
	Player     *Player
	RoundUUID  string
	Round      *Round `json:"-"`
	Value      int64
	Timestamp  time.Time
}
