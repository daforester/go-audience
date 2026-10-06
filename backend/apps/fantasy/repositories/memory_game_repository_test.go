package repositories

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
)

func TestConcurrentAccessAndSnapshot(t *testing.T) {
	repo := (&MemoryGameRepository{}).New()
	file := filepath.Join(t.TempDir(), "games.dat")

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			game := models.NewGame()
			game.Key = fmt.Sprintf("G%04d", i)
			_ = repo.Save(game)

			for j := 0; j < 20; j++ {
				g := repo.GetByKey(game.Key)
				_ = repo.Lock(g)
				g.AddPlayer(&models.Player{UUID: fmt.Sprintf("%d-%d", i, j), Name: "p"})
				_ = repo.Unlock(g)
				if j%5 == 0 {
					_ = repo.SaveToDisk(file)
				}
			}
		}(i)
	}
	wg.Wait()

	if err := repo.SaveToDisk(file); err != nil {
		t.Fatal(err)
	}

	loaded := (&MemoryGameRepository{}).New()
	if err := loaded.LoadFromDisk(file); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		g := loaded.GetByKey(fmt.Sprintf("G%04d", i))
		if g == nil || len(g.Players) != 20 {
			t.Fatalf("game %d not restored correctly", i)
		}
		if err := loaded.Lock(g); err != nil {
			t.Fatal(err)
		}
		_ = loaded.Unlock(g)
	}
}

func TestUnlockAfterDelete(t *testing.T) {
	repo := (&MemoryGameRepository{}).New()
	game := models.NewGame()
	game.Key = "ABCDE"
	_ = repo.Save(game)

	_ = repo.Lock(game)
	_ = repo.Delete(game)
	if err := repo.Unlock(game); err != nil {
		t.Fatalf("a deleted game must still unlock: %v", err)
	}
}
