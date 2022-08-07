package blogposts_test

import (
	"errors"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"

	blogposts "github.com/jrvldam/learn-go-with-tests/reading-files"
)

func TestNewBlogPosts(t *testing.T) {
	const (
		firstBody = `Title: Post 1
Description: Description 1
Tags: tdd, go`
		secondBody = `Title: Post 2
Description: Description 2
Tags: rust, borrow-checker`
	)
	fs := fstest.MapFS{
		"hello world.md": {Data: []byte(firstBody)},
		"hello-world.md": {Data: []byte(secondBody)},
	}

	posts, err := blogposts.NewPostsFromFS(fs)
	if err != nil {
		t.Fatal(err)
	}

	got := posts[0]
	want := blogposts.Post{
		Title:       "Post 1",
		Description: "Description 1",
		Tags:        []string{"tdd", "go"},
	}

	assertPost(t, got, want)
}

func TestNewBlogPostsFail(t *testing.T) {
	stub := StubFailingFs{}

	_, err := blogposts.NewPostsFromFS(stub)

	if err == nil {
		t.Errorf("expected to fail on reading a file")
	}
}

type StubFailingFs struct{}

func (s StubFailingFs) Open(name string) (fs.File, error) {
	return nil, errors.New("oh no. I always fail")
}

func assertPost(t *testing.T, got blogposts.Post, want blogposts.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
