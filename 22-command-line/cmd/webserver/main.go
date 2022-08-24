package main

import (
	"github.com/jrvldam/learn-go-with-tests/22-command-line"
	"log"
	"net/http"
	"os"
)

const dbFilename = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFilename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem opening %s %v", dbFilename, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)
	if err != nil {
		log.Fatalf("problem creating file system player store, %v", err)
	}

	server := poker.NewPlayerServer(store)

	err = http.ListenAndServe(":3000", server)
	if err != nil {
		log.Fatalf("could not listen on port 3000 %v", err)
	}
}
