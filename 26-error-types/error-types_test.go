package errortypes

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type BadStatusError struct {
	URL    string
	Status int
}

func (b BadStatusError) Error() string {
	return fmt.Sprintf("did not get 200 from %s, got %d", b.URL, b.Status)
}

func DumbGetter(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("problem fetching from %s, %v", url, err)
	}

	if res.StatusCode != http.StatusOK {
		return "", BadStatusError{URL: url, Status: res.StatusCode}
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	return string(body), nil
}

func TestDumbGetter(t *testing.T) {
	t.Run("when you do not get a 200 you get a status error", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
			res.WriteHeader(http.StatusTeapot)
		}))
		defer svr.Close()

		_, err := DumbGetter(svr.URL)

		if err == nil {
			t.Fatal("expected an error")
		}

		var got BadStatusError
		isBadStatusError := errors.As(err, &got)
		want := BadStatusError{URL: svr.URL, Status: http.StatusTeapot}

		if !isBadStatusError {
			t.Fatalf("was not BadStatusError, got %T", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}