package poker_test

import (
	"strings"
	"testing"

	poker "github.com/jrvldam/learn-go-with-tests/22-command-line"
)

func TestCLI(t *testing.T) {
	t.Run("record amaya win from user input", func(t *testing.T) {
		in := strings.NewReader("Amaya wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Amaya")
	})

	t.Run("record nayra win from user input", func(t *testing.T) {
		in := strings.NewReader("Nayra wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Nayra")
	})
}
