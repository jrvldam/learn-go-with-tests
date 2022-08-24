package main

import (
	"fmt"
	"log"
	"os"

	poker "github.com/jrvldam/learn-go-with-tests/22-command-line"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	pwd, err := os.Getwd()
	fmt.Println(":::::::", pwd)
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	poker.NewCLI(store, os.Stdin).PlayPoker()
}
