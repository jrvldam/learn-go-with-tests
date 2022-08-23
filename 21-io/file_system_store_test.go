package main

import (
	"strings"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		database := strings.NewReader(`
      [{"Name": "Nayra", "Wins": 10},
      {"Name": "Amaya", "Wins": 33}]
    `)
		store := FileSystemPlayerStore{database}

		got := store.GetLeague()
		want := []Player{
			{"Nayra", 10},
			{"Amaya", 33},
		}

		assertLeague(t, got, want)

		got = store.GetLeague()
		assertLeague(t, got, want)
	})
}
