package repositories

import (
	"bitbucket.org/daforester/go-audience/apps/trudge/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type lockableGame struct {
	Game *models.Game
	mux  sync.Mutex
}

type MemoryGameRepository struct {
	games map[string]*lockableGame
	mux   sync.Mutex
}

func (r MemoryGameRepository) New() *MemoryGameRepository {
	repo := new(MemoryGameRepository)
	repo.games = make(map[string]*lockableGame)
	return repo
}

func (r *MemoryGameRepository) Delete(game *models.Game) error {
	if game == nil {
		return errors.New("nil pointer passed")
	}

	_, exists := r.games[game.Key]
	if !exists {
		return errors.New("game does not exist")
	}
	delete(r.games, game.Key)
	return nil
}

func (r *MemoryGameRepository) GetByKey(key string) *models.Game {
	g, exists := r.games[key]
	if !exists {
		return nil
	}
	return g.Game
}

func (r *MemoryGameRepository) Save(game *models.Game) error {
	if game == nil {
		return errors.New("nil pointer passed")
	}

	g, exists := r.games[game.Key]
	if exists {
		g.Game = game
	} else {
		g = new(lockableGame)
		g.Game = game
		r.games[game.Key] = g
	}
	return nil
}

func (r *MemoryGameRepository) Lock(game *models.Game) error {
	if game == nil {
		r.mux.Lock()
		return nil
	}

	g, exists := r.games[game.Key]
	if !exists {
		return errors.New("game does not exist")
	}
	g.mux.Lock()
	return nil
}

func (r *MemoryGameRepository) Unlock(game *models.Game) error {
	if game == nil {
		r.mux.Unlock()
		return nil
	}

	g, exists := r.games[game.Key]
	if !exists {
		return errors.New("game does not exist")
	}
	g.mux.Unlock()
	return nil
}

func (r *MemoryGameRepository) LoadFromDisk(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = r.Lock(nil)

	if err != nil {
		return err
	}

	uErr := json.Unmarshal(data, &r.games)

	err = r.Unlock(nil)
	if err != nil {
		return err
	}

	return uErr
}

func (r *MemoryGameRepository) SaveToDisk(file string) error {
	var wErr error

	err := r.Lock(nil)

	if err != nil {
		return err
	}

	data, jErr := json.Marshal(r.games)
	fmt.Println(string(data))

	if jErr == nil {
		wErr = ioutil.WriteFile(file, data, os.ModePerm)
	}

	err = r.Unlock(nil)
	if err != nil {
		return err
	}

	if jErr != nil {
		return jErr
	}

	return wErr
}
