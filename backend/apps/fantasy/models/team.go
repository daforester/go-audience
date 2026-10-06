package models

type Team struct {
	ID         string
	Name       string
	RoundWins  int
	TotalVotes int
}

func (T *Team) Copy() *Team {
	t := *T
	return &t
}
