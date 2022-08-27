package poker_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	poker "github.com/jrvldam/learn-go-with-tests/24-websockets"
)

func TestGame_Start(t *testing.T) {
	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &poker.SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		game.Start(5)

		cases := []poker.ScheduleAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		checkSchedulingCases(t, cases, blindAlerter)
	})

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		stdin := strings.NewReader("7\n")
		blindAlerter := &poker.SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		cli := poker.NewCLI(stdin, stdout, game)
		cli.PlayPoker()

		got := poker.PlayerPrompt
		want := "Please enter the number of players: "

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		cases := []poker.ScheduleAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		checkSchedulingCases(t, cases, blindAlerter)
	})
}

func TestGame_Finish(t *testing.T) {
	store := &poker.StubPlayerStore{}
	game := poker.NewGame(dummyBlindAlerter, store)
	winner := "Julia"

	game.Finish(winner)

	poker.AssertPlayerWin(t, store, winner)
}

func checkSchedulingCases(t *testing.T, cases []poker.ScheduleAlert, blindAlerter *poker.SpyBlindAlerter) {
	for i, want := range cases {
		t.Run(fmt.Sprintf("%d scheduled for %v", want.Amount, want.At), func(t *testing.T) {
			if len(blindAlerter.Alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.Alerts)
			}

			got := blindAlerter.Alerts[i]
			assertScheduleAlert(t, got, want)
		})
	}
}

func assertScheduleAlert(t testing.TB, got poker.ScheduleAlert, want poker.ScheduleAlert) {
	t.Helper()
	amountGot := got.Amount
	if amountGot != want.Amount {
		t.Fatalf("got amount %d, want %d", amountGot, want.At)
	}

	gotScheduledTime := got.At
	if gotScheduledTime != want.At {
		t.Errorf("got scheduled time of %v, want %v", gotScheduledTime, want.At)
	}
}
