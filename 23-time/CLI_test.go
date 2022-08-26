package poker_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	poker "github.com/jrvldam/learn-go-with-tests/23-time"
)

var dummyBlindAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

type scheduleAlert struct {
	at     time.Duration
	amount int
}

func (s scheduleAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

type SpyBlindAlerter struct {
	alerts []scheduleAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduleAlert{at, amount})
}

func TestCLI(t *testing.T) {
	t.Run("record amaya win from user input", func(t *testing.T) {
		in := strings.NewReader("1\nAmaya wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in, dummyStdOut, dummyBlindAlerter)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Amaya")
	})

	t.Run("record nayra win from user input", func(t *testing.T) {
		in := strings.NewReader("1\nNayra wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in, dummyStdOut, dummyBlindAlerter)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Nayra")
	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		stdin := strings.NewReader("5\n")
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		cli := poker.NewCLI(playerStore, stdin, stdout, blindAlerter)
		cli.PlayPoker()

		cases := []scheduleAlert{
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

		for i, want := range cases {
			t.Run(fmt.Sprintf("%d scheduled for %v", want.amount, want.at), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduleAlert(t, got, want)
			})
		}
	})

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		stdin := strings.NewReader("7\n")
		blindAlerter := &SpyBlindAlerter{}

		cli := poker.NewCLI(dummyPlayerStore, stdin, stdout, blindAlerter)
		cli.PlayPoker()

		got := poker.PlayerPrompt
		want := "Please enter the number of players: "

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		cases := []scheduleAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		for i, want := range cases {
			t.Run(fmt.Sprintf("%d scheduled for %v", want.amount, want.at), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduleAlert(t, got, want)
			})
		}
	})
}

func assertScheduleAlert(t testing.TB, got scheduleAlert, want scheduleAlert) {
	t.Helper()
	amountGot := got.amount
	if amountGot != want.amount {
		t.Fatalf("got amount %d, want %d", amountGot, want.at)
	}

	gotScheduledTime := got.at
	if gotScheduledTime != want.at {
		t.Errorf("got scheduled time of %v, want %v", gotScheduledTime, want.at)
	}
}
