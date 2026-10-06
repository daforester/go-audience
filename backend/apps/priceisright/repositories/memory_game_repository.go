package repositories

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type MemoryGameRepository struct {
	game *models.Game
	mux  sync.Mutex
}

func (MemoryGameRepository) New(game *models.Game) *MemoryGameRepository {
	repo := new(MemoryGameRepository)
	if game != nil {
		repo.game = game
	} else {
		repo.game = models.Game{}.New()
	}
	return repo
}

func (r *MemoryGameRepository) Get() *models.Game {
	if r.game == nil {
		r.game = models.Game{}.New()
	}
	return r.game
}

func (r *MemoryGameRepository) Save(g *models.Game) {
	r.game = g
}

func (r *MemoryGameRepository) Lock() {
	r.mux.Lock()
}

func (r *MemoryGameRepository) Unlock() {
	r.mux.Unlock()
}

func (r *MemoryGameRepository) LoadFromDisk(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	r.Lock()

	uErr := json.Unmarshal(data, &r.game)
	if r.game.RoundNumber < 0 || r.game.RoundNumber > len(r.game.Rounds)-1 {
		r.game.RoundNumber = 0
	}
	if r.game.RoundNumber < len(r.game.Rounds) {
		r.game.CurrentRound = r.game.Rounds[r.game.RoundNumber]
	} else {
		r.game.CurrentRound = nil
	}

	r.Unlock()

	return uErr
}

func (r *MemoryGameRepository) SaveToDisk(file string) error {
	var wErr error

	r.Lock()

	data, jErr := json.Marshal(r.game)
	fmt.Println(string(data))

	if jErr == nil {
		wErr = ioutil.WriteFile(file, data, os.ModePerm)
	}

	r.Unlock()

	if jErr != nil {
		return jErr
	}

	return wErr
}
