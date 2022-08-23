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

	t.Run("get player score", func(t *testing.T) {
		database := strings.NewReader(`
      [{"Name": "Nayra", "Wins": 10},
      {"Name": "Amaya", "Wins": 33}]
    `)
		store := FileSystemPlayerStore{database}

		got := store.GetPlayerScore("Amaya")

		assertScoreEquals(t, got, 33)
	})
}

func assertScoreEquals(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
