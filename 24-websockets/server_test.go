package poker

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
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
	server := mustMakePlayerServer(t, &store)

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
	server := mustMakePlayerServer(t, &store)

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
		server := mustMakePlayerServer(t, &store)

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
		server := mustMakePlayerServer(t, &StubPlayerStore{})

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
	})

	t.Run("when we get a message over a websocket it is a winner of a game", func(t *testing.T) {
		store := &StubPlayerStore{}
		winner := "Amaya"
		server := httptest.NewServer(mustMakePlayerServer(t, store))
		defer server.Close()

		wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
		ws := mustDialWs(t, wsUrl)
		defer ws.Close()

		writeWsMessage(t, ws, winner)

		time.Sleep(10 * time.Millisecond)
		AssertPlayerWin(t, store, winner)
	})
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

func mustMakePlayerServer(t *testing.T, store PlayerStore) *PlayerServer {
	server, err := NewPlayerServer(store)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}

	return server
}
