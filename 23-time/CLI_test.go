package poker_test

import (
	"bytes"
	"strings"
	"testing"

	poker "github.com/jrvldam/learn-go-with-tests/23-time"
)

var dummyBlindAlerter = &poker.SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

type GameSpy struct {
	StartCalledWith  int
	FinishCalledWith string
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartCalledWith = numberOfPlayers
}

func (g *GameSpy) Finish(winner string) {
	g.FinishCalledWith = winner
}

func TestCLI(t *testing.T) {
	t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		stdin := strings.NewReader("7\n")
		game := &GameSpy{}

		cli := poker.NewCLI(stdin, stdout, game)
		cli.PlayPoker()

		gotPrompt := stdout.String()
		wantPrompt := poker.PlayerPrompt

		if gotPrompt != wantPrompt {
			t.Errorf("got %q, want %q", gotPrompt, wantPrompt)
		}

		if game.StartCalledWith != 7 {
			t.Errorf("wantend Start called with 7 but got %d", game.StartCalledWith)
		}
	})

	t.Run("record amaya win from user input", func(t *testing.T) {
		in := strings.NewReader("1\nAmaya wins\n")
		playerStore := &poker.StubPlayerStore{}
		game := poker.NewGame(dummyBlindAlerter, playerStore)

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Amaya")
	})

	t.Run("record nayra win from user input", func(t *testing.T) {
		in := strings.NewReader("1\nNayra wins\n")
		playerStore := &poker.StubPlayerStore{}
		game := poker.NewGame(dummyBlindAlerter, playerStore)

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Nayra")
	})
}
