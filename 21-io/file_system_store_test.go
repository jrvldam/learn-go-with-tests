package main

import (
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `
      [{"Name": "Nayra", "Wins": 10},
      {"Name": "Amaya", "Wins": 33}]
    `)
		defer cleanDatabase()
		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		got := store.GetLeague()
		want := []Player{
			{"Amaya", 33},
			{"Nayra", 10},
		}

		assertLeague(t, got, want)

		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `
      [{"Name": "Nayra", "Wins": 10},
      {"Name": "Amaya", "Wins": 33}]
    `)
		defer cleanDatabase()
		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		got := store.GetPlayerScore("Amaya")

		assertScoreEquals(t, got, 33)
	})

	t.Run("store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `
      [{"Name": "Nayra", "Wins": 10},
      {"Name": "Amaya", "Wins": 33}]
     `)
		defer cleanDatabase()
		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		store.RecordWin("Amaya")
		got := store.GetPlayerScore("Amaya")

		assertScoreEquals(t, got, 34)
	})

	t.Run("store wins for new players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `
      [{"Name": "Nayra", "Wins": 10},
      {"Name": "Amaya", "Wins": 33}]
     `)
		defer cleanDatabase()
		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		store.RecordWin("Julia")
		got := store.GetPlayerScore("Julia")

		assertScoreEquals(t, got, 1)

	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)
	})

	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `
      [{"Name": "Nayra", "Wins": 10},
      {"Name": "Julia", "Wins": 33}]
    `)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		got := store.GetLeague()
		want := []Player{
			{"Julia", 33},
			{"Nayra", 10},
		}

		assertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})
}

func assertScoreEquals(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "db")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpFile.Write([]byte(initialData))

	removeFile := func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	return tmpFile, removeFile
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("did not expected an error but got one, %v", err)
	}
}
