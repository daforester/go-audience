package models

import (
	"bitbucket.org/daforester/go-audience/libs/slice"
	"sort"
	"time"
)

type Round struct {
	Order     int
	Image     string
	Name      string
	Value     int64
	Guesses   []*Guess
	UUID      string
	Activated time.Time
}

func (r *Round) AddGuess(g *Guess) {
	r.Guesses = append(r.Guesses, g)
}

func (r *Round) TopGuesses() []*Guess {
	g := make([]*Guess, len(r.Guesses))
	copy(g, r.Guesses)

	// First remove any values over the actual value
	slice.Filter(&g, func(v any) bool {
		return v.(*Guess).Value <= r.Value
	})

	// If there aren't 2 remaining results, there's going to be a problem
	if len(g) < 2 {
		return nil
	}

	// Now order the guesses by value, high to low
	sort.Slice(g, func(i, j int) bool {
		return g[i].Value > r.Guesses[j].Value
	})

	// Get the top two closest values - they might be the same
	v1 := g[0].Value
	v2 := g[1].Value

	// Remove anything that's not a top 2 guess
	slice.Filter(&g, func(v any) bool {
		return v.(*Guess).Value == v1 || v.(*Guess).Value == v2
	})

	// Order remaining results by time of submission
	sort.Slice(g, func(i, j int) bool {
		return g[i].Timestamp.UnixNano() < g[j].Timestamp.UnixNano()
	})

	if len(g) > 10 {
		return g[0:9]
	}
	return g
}
