package poker

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var (
	dummyStore = &StubPlayerStore{}
	dummyGame  = &SpyGame{}
)

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := mustMakePlayerServer(t, &store, dummyGame)

	tests := []struct {
		name               string
		player             string
		expectedHttpStatus int
		expectedScore      string
	}{
		{
			name:               "Returns Pepper's score",
			player:             "Pepper",
			expectedHttpStatus: http.StatusOK,
			expectedScore:      "20",
		},
		{
			name:               "Returns Floyd's score",
			player:             "Floyd",
			expectedHttpStatus: http.StatusOK,
			expectedScore:      "10",
		},
		{
			name:               "Returns 404 on missing players",
			player:             "Appolo",
			expectedHttpStatus: http.StatusNotFound,
			expectedScore:      "0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := newGetScoreRequest(test.player)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertStatus(t, response, test.expectedHttpStatus)
			assertResponseBody(t, response.Body.String(), test.expectedScore)
		})
	}
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := mustMakePlayerServer(t, &store, dummyGame)

	t.Run("it records on POST", func(t *testing.T) {
		player := "Pepper"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusAccepted)
		AssertPlayerWin(t, &store, player)
	})
}

func TestLeague(t *testing.T) {
	t.Run("it returns the league table as json", func(t *testing.T) {
		wantedLeague := []Player{
			{"Nayra", 32},
			{"Julia", 20},
			{"Amaya", 14},
		}

		store := StubPlayerStore{nil, nil, wantedLeague}
		server := mustMakePlayerServer(t, &store, dummyGame)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)

		assertStatus(t, response, http.StatusOK)
		assertLeague(t, got, wantedLeague)
		assertContentType(t, response, jsonContentType)
	})
}

func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := mustMakePlayerServer(t, &StubPlayerStore{}, dummyGame)

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
	})

	t.Run("start a game with 3 players and declare Amaya the winner", func(t *testing.T) {
		game := &SpyGame{}
		winner := "Amaya"
		server := httptest.NewServer(mustMakePlayerServer(t, dummyStore, game))
		ws := mustDialWs(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")
		defer server.Close()
		defer ws.Close()

		writeWsMessage(t, ws, "3")
		writeWsMessage(t, ws, winner)

		time.Sleep(10 * time.Millisecond)
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, winner)
	})
}

func assertGameStartedWith(t testing.TB, game *SpyGame, players int) {
	t.Helper()

	if game.StartCalledWith != players {
		t.Errorf("got %d, want %d", game.StartCalledWith, players)
	}
}

func assertFinishCalledWith(t testing.TB, game *SpyGame, winner string) {
	t.Helper()

	if game.FinishCalledWith != winner {
		t.Errorf("got %s, want %s", game.FinishCalledWith, winner)
	}
}

func writeWsMessage(t testing.TB, conn *websocket.Conn, message string) {
	t.Helper()

	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}

func mustDialWs(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s, %v", url, err)
	}

	return ws
}

func mustMakePlayerServer(t *testing.T, store PlayerStore, game *SpyGame) *PlayerServer {
	server, err := NewPlayerServer(store, game)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}

	return server
}
