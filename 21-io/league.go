package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func NewLeague(rdr io.Reader) (League, error) {
	var league League
	err := json.NewDecoder(rdr).Decode(&league)

	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}

	return league, err
}

type League []Player

func (l League) Find(name string) *Player {
	for idx, player := range l {
		if player.Name == name {
			return &l[idx]
		}
	}

	return nil
}
