package main

import (
	"log"
	"net/http"

	poker "github.com/jrvldam/learn-go-with-tests/24-websockets"
)

const dbFilename = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	server := poker.NewPlayerServer(store)

	err = http.ListenAndServe(":3000", server)
	if err != nil {
		log.Fatalf("could not listen on port 3000 %v", err)
	}
}
