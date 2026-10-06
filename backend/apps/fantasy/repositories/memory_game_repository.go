package repositories

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type lockableGame struct {
	Game *models.Game
	mux  *sync.Mutex
}

// MemoryGameRepository keeps games in memory.
//
//   - mapMux guards the games map itself, every method takes it.
//   - Lock(nil) / Unlock(nil) take mux, used for whole-repository operations such as
//     reserving a unique room key.
//   - Lock(game) / Unlock(game) take that game's own mutex. Processors hold it while they
//     read or change the game. Game mutexes outlive Delete, so a processor that deletes a
//     game can still Unlock it.
type MemoryGameRepository struct {
	games  map[string]*lockableGame
	locks  map[string]*sync.Mutex
	mapMux sync.RWMutex
	mux    sync.Mutex
}

func (r *MemoryGameRepository) New() *MemoryGameRepository {
	repo := new(MemoryGameRepository)
	repo.games = make(map[string]*lockableGame)
	repo.locks = make(map[string]*sync.Mutex)
	return repo
}

func (r *MemoryGameRepository) Delete(game *models.Game) error {
	if game == nil {
		return errors.New("nil pointer passed")
	}

	r.mapMux.Lock()
	defer r.mapMux.Unlock()

	_, exists := r.games[game.Key]
	if !exists {
		return errors.New("game does not exist")
	}
	delete(r.games, game.Key)
	return nil
}

func (r *MemoryGameRepository) GetByKey(key string) *models.Game {
	r.mapMux.RLock()
	defer r.mapMux.RUnlock()

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

	r.mapMux.Lock()
	defer r.mapMux.Unlock()

	g, exists := r.games[game.Key]
	if exists {
		g.Game = game
	} else {
		g = new(lockableGame)
		g.Game = game
		g.mux = r.gameMutex(game.Key)
		r.games[game.Key] = g
	}
	return nil
}

func (r *MemoryGameRepository) Lock(game *models.Game) error {
	if game == nil {
		r.mux.Lock()
		return nil
	}

	m := r.existingMutex(game.Key)
	if m == nil {
		return errors.New("game does not exist")
	}
	m.Lock()
	return nil
}

func (r *MemoryGameRepository) Unlock(game *models.Game) error {
	if game == nil {
		r.mux.Unlock()
		return nil
	}

	m := r.existingMutex(game.Key)
	if m == nil {
		return errors.New("game does not exist")
	}
	m.Unlock()
	return nil
}

func (r *MemoryGameRepository) existingMutex(key string) *sync.Mutex {
	r.mapMux.RLock()
	defer r.mapMux.RUnlock()
	return r.locks[key]
}

// gameMutex must be called with mapMux held for writing
func (r *MemoryGameRepository) gameMutex(key string) *sync.Mutex {
	m, exists := r.locks[key]
	if !exists {
		m = new(sync.Mutex)
		r.locks[key] = m
	}
	return m
}

func (r *MemoryGameRepository) LoadFromDisk(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	games := make(map[string]*lockableGame)
	err = json.Unmarshal(data, &games)
	if err != nil {
		return err
	}

	r.mapMux.Lock()
	defer r.mapMux.Unlock()

	for key, g := range games {
		if g == nil || g.Game == nil {
			continue
		}
		g.mux = r.gameMutex(key)
		r.games[key] = g
	}

	return nil
}

func (r *MemoryGameRepository) SaveToDisk(file string) error {
	r.mapMux.RLock()
	games := make(map[string]*lockableGame, len(r.games))
	for k, g := range r.games {
		games[k] = g
	}
	r.mapMux.RUnlock()

	// Marshal each game while holding its own lock so we never read a game mid-change
	out := make(map[string]map[string]json.RawMessage, len(games))
	for k, g := range games {
		g.mux.Lock()
		data, err := json.Marshal(g.Game)
		g.mux.Unlock()
		if err != nil {
			return err
		}
		out[k] = map[string]json.RawMessage{"Game": data}
	}

	data, err := json.Marshal(out)
	if err != nil {
		return err
	}

	// Write to a temporary file and rename, so a crash mid-write never corrupts the state file
	tmp, err := os.CreateTemp(filepath.Dir(file), filepath.Base(file)+".*.tmp")
	if err != nil {
		return err
	}
	_, err = tmp.Write(data)
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(tmp.Name())
		return err
	}

	return os.Rename(tmp.Name(), file)
}
