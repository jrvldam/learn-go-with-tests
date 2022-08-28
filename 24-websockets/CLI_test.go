package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	poker "github.com/jrvldam/learn-go-with-tests/24-websockets"
)

var dummyBlindAlerter = &poker.SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}

func TestCLI(t *testing.T) {
	t.Run("start game with 3 players and finish game with 'Nayra' as winner", func(t *testing.T) {
		game := &poker.SpyGame{}
		stdout := &bytes.Buffer{}

		stdin := userSends("3", "Nayra wins")
		cli := poker.NewCLI(stdin, stdout, game)

		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, "Nayra")
	})

	t.Run("start game with 8 players and record 'Julia' as winner", func(t *testing.T) {
		game := &poker.SpyGame{}

		stdin := userSends("8", "Julia wins")
		cli := poker.NewCLI(stdin, dummyStdOut, game)

		cli.PlayPoker()

		assertGameStartedWith(t, game, 8)
		assertFinishCalledWith(t, game, "Julia")
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		stdin := strings.NewReader("Rodrigo\n")
		game := &poker.SpyGame{}

		cli := poker.NewCLI(stdin, stdout, game)
		cli.PlayPoker()

		gotPrompt := stdout.String()
		wantPrompt := poker.PlayerPrompt + poker.BadPlayerInputErrMsg

		if gotPrompt != wantPrompt {
			t.Errorf("got %q, want %q", gotPrompt, wantPrompt)
		}

		if game.StartCalled {
			t.Errorf("game should not have started")
		}
	})

	t.Run("it prints an error when the winner is declared incorrectly", func(t *testing.T) {
		game := &poker.SpyGame{}
		stdout := &bytes.Buffer{}

		stdin := userSends("1", "Rodrigo is a killer")
		cli := poker.NewCLI(stdin, stdout, game)

		cli.PlayPoker()

		assertGameNotFinished(t, game)
		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.BadWinnerInputMsg)
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

func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	got := stdout.String()
	want := strings.Join(messages, "")

	if got != want {
		t.Errorf("got %q sent to stdout, but expected %+v", got, messages)
	}
}

func assertGameStartedWith(t testing.TB, game *poker.SpyGame, want int) {
	t.Helper()

	got := game.StartCalledWith

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func assertFinishCalledWith(t testing.TB, game *poker.SpyGame, want string) {
	t.Helper()

	got := game.FinishCalledWith

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func assertGameNotFinished(t testing.TB, game *poker.SpyGame) {
	t.Helper()
	if game.FinishCalled {
		t.Errorf("finish should not to be called")
	}
}
