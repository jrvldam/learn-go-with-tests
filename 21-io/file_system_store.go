package main

import (
	"io"
)

type FileSystemPlayerStore struct {
	database io.ReadSeeker
}

func (f *FileSystemPlayerStore) GetLeague() []Player {
	// Back to the first position in order to re-read
	f.database.Seek(0, 0)

	league, _ := NewLeague(f.database)

	return league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	var score int
	for _, player := range f.GetLeague() {
		if player.Name == name {
			score = player.Wins
			return score
		}
	}

	return score
}
