package poker

import (
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	t.Run("record amaya win from user input", func(t *testing.T) {
		in := strings.NewReader("Amaya wins\n")
		playerStore := &StubPlayerStore{}

		cli := &CLI{playerStore, in}
		cli.PlayPoker()

		assertPlayerWin(t, playerStore, "Amaya")
	})

	t.Run("record nayra win from user input", func(t *testing.T) {
		in := strings.NewReader("Nayra wins\n")
		playerStore := &StubPlayerStore{}

		cli := &CLI{playerStore, in}
		cli.PlayPoker()

		assertPlayerWin(t, playerStore, "Nayra")
	})
}
